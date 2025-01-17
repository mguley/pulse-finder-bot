package application

import (
	"application/config"
	"application/dependency"
	"application/proxy/circuit"
	"application/proxy/commands"
	"application/proxy/commands/control"
	"application/proxy/services"
	"application/proxy/strategies"
	appScheduler "application/scheduler"
	"application/source"
	sourceAlfa "application/source/alfa"
	sourceBeta "application/source/beta"
	"application/url/processor"
	"application/url/sitemap"
	"domain/html"
	"domain/scheduler"
	"infrastructure"
	htmlAlfa "infrastructure/html/source/alfa"
	htmlBeta "infrastructure/html/source/beta"
	"infrastructure/url/sitemap/fetcher"
	"infrastructure/url/sitemap/notifier"
	"infrastructure/url/sitemap/parser"
	sitemapRepository "infrastructure/url/sitemap/repository"
	"log"
	"net/http"
	"time"
)

// Container is a struct that holds all the dependencies for the application.
// It acts as a central registry for services, ensuring that dependencies are managed in a lazy loaded manner.
type Container struct {
	Config                  dependency.LazyDependency[*config.Config]
	SitemapServiceXML       dependency.LazyDependency[*sitemap.Service]
	SitemapServiceRSS       dependency.LazyDependency[*sitemap.Service]
	SitemapFetcher          dependency.LazyDependency[*fetcher.Service]
	SitemapNotifier         dependency.LazyDependency[*notifier.Service]
	SitemapParser           dependency.LazyDependency[*parser.Service]
	SitemapParserRSS        dependency.LazyDependency[*parser.RssFeed]
	SitemapRepository       dependency.LazyDependency[*sitemapRepository.Service]
	ProxyService            dependency.LazyDependency[*services.Service]
	RetryStrategy           dependency.LazyDependency[strategies.RetryStrategy]
	IdentityService         dependency.LazyDependency[*services.Identity]
	CircuitManager          dependency.LazyDependency[*circuit.Manager]
	AlfaHtmlFetcher         dependency.LazyDependency[html.Fetcher]
	AlfaHtmlParser          dependency.LazyDependency[html.Parser]
	AlfaHandler             dependency.LazyDependency[*sourceAlfa.Handler]
	BetaHtmlFetcher         dependency.LazyDependency[html.Fetcher]
	BetaHtmlParser          dependency.LazyDependency[html.Parser]
	BetaHandler             dependency.LazyDependency[*sourceBeta.Handler]
	SourceFactory           dependency.LazyDependency[*source.Factory]
	ProcessorService        dependency.LazyDependency[*processor.Service]
	AuthenticateCommand     dependency.LazyDependency[*control.AuthenticateCommand]
	SignalCommand           dependency.LazyDependency[*control.SignalCommand]
	StatusCommand           dependency.LazyDependency[*commands.StatusCommand]
	InfrastructureContainer dependency.LazyDependency[*infrastructure.Container]
	CronScheduler           dependency.LazyDependency[scheduler.Scheduler]
}

// NewContainer creates and returns a new instance of Container.
// Each dependency is configured to initialize only when first accessed.
func NewContainer() *Container {
	c := &Container{}

	// Create container with base dependencies
	c.Config = dependency.LazyDependency[*config.Config]{
		InitFunc: config.LoadConfig,
	}

	// Proxy commands
	c.AuthenticateCommand = dependency.LazyDependency[*control.AuthenticateCommand]{
		InitFunc: func() *control.AuthenticateCommand {
			conn := c.InfrastructureContainer.Get().ProxyConnection.Get()
			return control.NewAuthenticateCommand(conn)
		},
	}
	c.SignalCommand = dependency.LazyDependency[*control.SignalCommand]{
		InitFunc: func() *control.SignalCommand {
			conn := c.InfrastructureContainer.Get().ProxyConnection.Get()
			return control.NewSignalCommand(conn, "NEWNYM")
		},
	}
	c.StatusCommand = dependency.LazyDependency[*commands.StatusCommand]{
		InitFunc: func() *commands.StatusCommand {
			factory := c.InfrastructureContainer.Get().HttpFactory.Get()
			client := c.InfrastructureContainer.Get().Socks5Client.Get()
			host := c.Config.Get().Proxy.Host
			port := c.Config.Get().Proxy.Port
			timeout := 10 * time.Second
			return commands.NewStatusCommand(host, port, client, factory, timeout)
		},
	}

	// Proxy services
	c.ProxyService = dependency.LazyDependency[*services.Service]{
		InitFunc: func() *services.Service {
			factory := c.InfrastructureContainer.Get().HttpFactory.Get()
			host := c.Config.Get().Proxy.Host
			port := c.Config.Get().Proxy.Port
			timeout := 10 * time.Second
			return services.NewService(factory, host, port, timeout)
		},
	}
	c.RetryStrategy = dependency.LazyDependency[strategies.RetryStrategy]{
		InitFunc: func() strategies.RetryStrategy {
			baseDelay := 5 * time.Second
			maxDelay := 30 * time.Second
			maxAttempts := 5
			multiplier := 2.0
			return strategies.NewExponentialBackoffStrategy(baseDelay, maxDelay, maxAttempts, multiplier)
		},
	}
	c.IdentityService = dependency.LazyDependency[*services.Identity]{
		InitFunc: func() *services.Identity {
			conn := c.InfrastructureContainer.Get().ProxyConnection.Get()
			auth := c.AuthenticateCommand.Get()
			signal := c.SignalCommand.Get()
			status := c.StatusCommand.Get()
			strategy := c.RetryStrategy.Get()
			url := c.Config.Get().Proxy.PingUrl
			return services.NewIdentity(conn, auth, signal, status, strategy, url)
		},
	}
	c.CircuitManager = dependency.LazyDependency[*circuit.Manager]{
		InitFunc: func() *circuit.Manager {
			return circuit.NewManager(c.IdentityService.Get(), c.ProxyService.Get(), c.Config.Get().Proxy.PingUrl)
		},
	}

	// Parser services
	c.AlfaHtmlFetcher = dependency.LazyDependency[html.Fetcher]{
		InitFunc: func() html.Fetcher {
			maxBodySize := int64(10 * 1024 * 1024) // 10MB
			htmlFetcher, err := htmlAlfa.NewFetcher(c.ProxyService.Get(), maxBodySize)
			if err != nil {
				log.Fatalf("failed to init fetcher: %v", err)
			}
			return htmlFetcher
		},
	}
	c.AlfaHtmlParser = dependency.LazyDependency[html.Parser]{
		InitFunc: func() html.Parser {
			return htmlAlfa.NewParser()
		},
	}
	c.AlfaHandler = dependency.LazyDependency[*sourceAlfa.Handler]{
		InitFunc: func() *sourceAlfa.Handler {
			url := c.Config.Get().SourceHandler.Alfa.SitemapURL
			sitemapService := c.SitemapServiceRSS.Get()
			circuitManager := c.CircuitManager.Get()
			urlRepository := c.InfrastructureContainer.Get().UrlRepository.Get()
			vacancyRepository := c.InfrastructureContainer.Get().VacancyRepository.Get()
			htmlFetcher := c.AlfaHtmlFetcher.Get()
			htmlParser := c.AlfaHtmlParser.Get()
			return sourceAlfa.NewHandler(url, sitemapService, circuitManager, urlRepository,
				vacancyRepository, htmlFetcher, htmlParser)
		},
	}
	c.BetaHtmlFetcher = dependency.LazyDependency[html.Fetcher]{
		InitFunc: func() html.Fetcher {
			maxBodySize := int64(10 * 1024 * 1024) // 10MB
			htmlFetcher, err := htmlBeta.NewFetcher(c.ProxyService.Get(), maxBodySize)
			if err != nil {
				log.Fatalf("failed to init fetcher: %v", err)
			}
			return htmlFetcher
		},
	}
	c.BetaHtmlParser = dependency.LazyDependency[html.Parser]{
		InitFunc: func() html.Parser {
			return htmlBeta.NewParser()
		},
	}
	c.BetaHandler = dependency.LazyDependency[*sourceBeta.Handler]{
		InitFunc: func() *sourceBeta.Handler {
			url := c.Config.Get().SourceHandler.Beta.SitemapURL
			sitemapService := c.SitemapServiceXML.Get()
			circuitManager := c.CircuitManager.Get()
			urlRepository := c.InfrastructureContainer.Get().UrlRepository.Get()
			vacancyRepository := c.InfrastructureContainer.Get().VacancyRepository.Get()
			htmlFetcher := c.BetaHtmlFetcher.Get()
			htmlParser := c.BetaHtmlParser.Get()
			return sourceBeta.NewHandler(url, sitemapService, circuitManager, urlRepository,
				vacancyRepository, htmlFetcher, htmlParser)
		},
	}

	// Domain/layer containers
	c.InfrastructureContainer = dependency.LazyDependency[*infrastructure.Container]{
		InitFunc: func() *infrastructure.Container {
			return infrastructure.NewContainer(c.Config.Get())
		},
	}

	// Sitemap services
	c.SitemapFetcher = dependency.LazyDependency[*fetcher.Service]{
		InitFunc: func() *fetcher.Service {
			proxyClient := func() (*http.Client, error) {
				return c.ProxyService.Get().HttpClient()
			}
			return fetcher.NewService(proxyClient)
		},
	}
	c.SitemapNotifier = dependency.LazyDependency[*notifier.Service]{
		InitFunc: func() *notifier.Service {
			return notifier.NewService(c.Config.Get().Proxy.PingUrl)
		},
	}
	c.SitemapParser = dependency.LazyDependency[*parser.Service]{
		InitFunc: parser.NewService,
	}
	c.SitemapParserRSS = dependency.LazyDependency[*parser.RssFeed]{
		InitFunc: parser.NewRssFeed,
	}
	c.SitemapRepository = dependency.LazyDependency[*sitemapRepository.Service]{
		InitFunc: func() *sitemapRepository.Service {
			return sitemapRepository.NewService(c.InfrastructureContainer.Get().UrlRepository.Get())
		},
	}
	c.SitemapServiceXML = dependency.LazyDependency[*sitemap.Service]{
		InitFunc: func() *sitemap.Service {
			return sitemap.NewService(
				sitemap.WithFetcher(c.SitemapFetcher.Get()),
				sitemap.WithParser(c.SitemapParser.Get()),
				sitemap.WithRepository(c.SitemapRepository.Get()),
				sitemap.WithNotifier(c.SitemapNotifier.Get()),
				sitemap.WithHTTPClient(func() (*http.Client, error) {
					return c.ProxyService.Get().HttpClient()
				}))
		},
	}
	c.SitemapServiceRSS = dependency.LazyDependency[*sitemap.Service]{
		InitFunc: func() *sitemap.Service {
			return sitemap.NewService(
				sitemap.WithFetcher(c.SitemapFetcher.Get()),
				sitemap.WithParser(c.SitemapParserRSS.Get()),
				sitemap.WithRepository(c.SitemapRepository.Get()),
				sitemap.WithNotifier(c.SitemapNotifier.Get()),
				sitemap.WithHTTPClient(func() (*http.Client, error) {
					return c.ProxyService.Get().HttpClient()
				}))
		},
	}
	c.SourceFactory = dependency.LazyDependency[*source.Factory]{
		InitFunc: source.NewFactory,
	}
	c.ProcessorService = dependency.LazyDependency[*processor.Service]{
		InitFunc: func() *processor.Service {
			sourceFactory := c.SourceFactory.Get()
			batchSize := c.Config.Get().SourceHandler.BatchSize
			return processor.NewService(sourceFactory, batchSize)
		},
	}

	// Cron
	c.CronScheduler = dependency.LazyDependency[scheduler.Scheduler]{
		InitFunc: func() scheduler.Scheduler {
			infra := c.InfrastructureContainer.Get()
			cfg := c.Config.Get()
			repo := infra.VacancyRepository.Get()
			aClient := infra.AuthClient.Get()
			vClient := infra.VacancyClient.Get()
			batchSize := 5
			issuer := cfg.AuthServer.Issuer
			scope := []string{"write"}
			tickerTime := 15 * time.Second
			return appScheduler.NewCronScheduler(repo, aClient, vClient, batchSize, issuer, scope, tickerTime)
		},
	}

	return c
}
