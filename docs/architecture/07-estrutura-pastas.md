# 7. Nova Estrutura de Pastas

## 🎯 Objetivo

Reorganizar a estrutura de pastas seguindo princípios de Clean Architecture e Domain-Driven Design para melhorar organização, manutenibilidade e escalabilidade.

## 🏗️ Estrutura Proposta

```bash
girus-cli/
├── cmd/                              # Application Layer - CLI Commands
│   ├── create.go                     # Apenas definição de comandos Cobra
│   ├── delete.go                     # Comandos de remoção
│   ├── list.go                       # Comandos de listagem
│   ├── lab.go                        # Comandos específicos de lab
│   ├── repo.go                       # Comandos de repositório
│   ├── status.go                     # Comandos de status
│   └── root.go                       # Comando raiz e configuração
│
├── internal/                         # Código interno da aplicação
│   │
│   ├── app/                          # Application Layer
│   │   ├── handlers/                 # Command Handlers
│   │   │   ├── cluster/              # Handlers de cluster
│   │   │   │   ├── create_handler.go
│   │   │   │   ├── delete_handler.go
│   │   │   │   ├── status_handler.go
│   │   │   │   └── interfaces.go
│   │   │   ├── lab/                  # Handlers de laboratório
│   │   │   │   ├── install_handler.go
│   │   │   │   ├── uninstall_handler.go
│   │   │   │   ├── list_handler.go
│   │   │   │   └── interfaces.go
│   │   │   └── shared/               # Handlers compartilhados
│   │   │       ├── base_handler.go
│   │   │       └── validation.go
│   │   │
│   │   ├── orchestrators/            # Application Services (Orchestration)
│   │   │   ├── cluster_orchestrator.go
│   │   │   ├── lab_orchestrator.go
│   │   │   └── infrastructure_orchestrator.go
│   │   │
│   │   ├── dto/                      # Data Transfer Objects
│   │   │   ├── cluster_dto.go
│   │   │   ├── lab_dto.go
│   │   │   └── common_dto.go
│   │   │
│   │   └── container.go              # Dependency Injection Container
│   │
│   ├── domain/                       # Domain Layer - Business Logic
│   │   ├── cluster/                  # Cluster Domain
│   │   │   ├── entities.go           # Cluster, Node, Status entities
│   │   │   ├── services.go           # Business logic interfaces
│   │   │   ├── repository.go         # Data access interfaces
│   │   │   ├── specifications.go     # Business rules
│   │   │   ├── events.go             # Domain events
│   │   │   └── errors.go             # Domain-specific errors
│   │   │
│   │   ├── lab/                      # Lab Domain
│   │   │   ├── entities.go           # Lab, Template, Task entities
│   │   │   ├── services.go           # Lab business logic
│   │   │   ├── repository.go         # Lab data access
│   │   │   ├── validation.go         # Lab validation rules
│   │   │   ├── specification.go      # Lab business rules
│   │   │   └── errors.go             # Lab-specific errors
│   │   │
│   │   ├── repository/               # Repository Domain
│   │   │   ├── entities.go           # Repository, Index entities
│   │   │   ├── services.go           # Repository operations
│   │   │   ├── repository.go         # Repository data access
│   │   │   └── errors.go             # Repository errors
│   │   │
│   │   ├── infrastructure/           # Infrastructure Domain
│   │   │   ├── entities.go           # System, Prerequisites entities
│   │   │   ├── services.go           # Infrastructure services
│   │   │   ├── detector.go           # System detection
│   │   │   └── installer.go          # Tool installation
│   │   │
│   │   └── shared/                   # Shared Domain Concepts
│   │       ├── value_objects.go      # Version, Resource, etc.
│   │       ├── events.go             # Domain events
│   │       ├── specifications.go     # Common business rules
│   │       └── interfaces.go         # Shared interfaces
│   │
│   ├── infrastructure/               # Infrastructure Layer
│   │   ├── k8s/                      # Kubernetes Infrastructure
│   │   │   ├── client.go             # K8s client wrapper
│   │   │   ├── cluster_repository.go # K8s cluster operations
│   │   │   ├── lab_repository.go     # ConfigMap operations
│   │   │   ├── deployment.go         # Deployment operations
│   │   │   └── port_forward.go       # Port forwarding
│   │   │
│   │   ├── containerengine/          # Container Engine Abstraction
│   │   │   ├── interfaces.go         # Container engine interfaces
│   │   │   ├── docker/               # Docker implementation
│   │   │   │   ├── client.go
│   │   │   │   ├── service.go
│   │   │   │   └── detector.go
│   │   │   ├── podman/               # Podman implementation
│   │   │   │   ├── client.go
│   │   │   │   ├── service.go
│   │   │   │   └── detector.go
│   │   │   └── factory.go            # Engine factory
│   │   │
│   │   ├── kind/                     # Kind-specific operations
│   │   │   ├── cluster_manager.go
│   │   │   ├── config_generator.go
│   │   │   └── node_manager.go
│   │   │
│   │   ├── http/                     # HTTP Clients
│   │   │   ├── repository_client.go  # Repository HTTP client
│   │   │   ├── version_client.go     # Version check client
│   │   │   ├── retry_client.go       # HTTP retry wrapper
│   │   │   └── auth_client.go        # Authentication
│   │   │
│   │   ├── filesystem/               # File System Operations
│   │   │   ├── file_repository.go    # File operations
│   │   │   ├── template_loader.go    # Template file loading
│   │   │   ├── cache_manager.go      # File caching
│   │   │   └── watcher.go            # File watching
│   │   │
│   │   ├── exec/                     # Command Execution
│   │   │   ├── command_executor.go   # Command execution wrapper
│   │   │   ├── process_manager.go    # Process management
│   │   │   ├── stream_handler.go     # Output streaming
│   │   │   └── timeout_manager.go    # Timeout handling
│   │   │
│   │   └── ui/                       # User Interface Infrastructure
│   │       ├── console/              # Console UI implementation
│   │       │   ├── console_ui.go
│   │       │   ├── progress_bars.go
│   │       │   ├── formatters.go
│   │       │   └── colors.go
│   │       ├── json/                 # JSON output
│   │       │   └── json_ui.go
│   │       └── interfaces.go         # UI interfaces
│   │
│   ├── adapters/                     # Interface Adapters Layer
│   │   ├── config/                   # Configuration Adapters
│   │   │   ├── loader.go             # Multi-source config loader
│   │   │   ├── validator.go          # Config validation
│   │   │   ├── merger.go             # Config merging logic
│   │   │   ├── types.go              # Config types
│   │   │   └── defaults.go           # Default values
│   │   │
│   │   ├── logging/                  # Logging Adapters
│   │   │   ├── structured_logger.go  # Structured logging
│   │   │   ├── file_logger.go        # File-based logging
│   │   │   ├── console_logger.go     # Console logging
│   │   │   └── interfaces.go         # Logger interfaces
│   │   │
│   │   ├── i18n/                     # Internationalization
│   │   │   ├── translator.go         # Translation service
│   │   │   ├── loader.go             # Message file loader
│   │   │   ├── formatter.go          # Message formatting
│   │   │   └── messages/             # Translation files
│   │   │       ├── en.yaml           # English messages
│   │   │       ├── pt.yaml           # Portuguese messages
│   │   │       └── es.yaml           # Spanish messages
│   │   │
│   │   ├── telemetry/                # Telemetry and Metrics
│   │   │   ├── metrics_collector.go  # Metrics collection
│   │   │   ├── error_reporter.go     # Error reporting
│   │   │   ├── usage_tracker.go      # Usage analytics
│   │   │   └── exporters/            # Various exporters
│   │   │       ├── prometheus.go
│   │   │       └── jaeger.go
│   │   │
│   │   └── errors/                   # Error Handling Adapters
│   │       ├── handler.go            # Main error handler
│   │       ├── recovery.go           # Recovery strategies
│   │       ├── formatters.go         # Error formatters
│   │       ├── types.go              # Error types
│   │       └── factories.go          # Error factories
│   │
│   ├── pkg/                          # Shared Internal Packages
│   │   ├── validation/               # Generic Validation
│   │   │   ├── rules.go              # Validation rules
│   │   │   ├── validator.go          # Validator implementation
│   │   │   └── custom_validators.go  # Custom validation functions
│   │   │
│   │   ├── retry/                    # Retry Mechanisms
│   │   │   ├── retry.go              # Retry logic
│   │   │   ├── backoff.go            # Backoff strategies
│   │   │   └── circuit_breaker.go    # Circuit breaker
│   │   │
│   │   ├── cache/                    # Caching Utilities
│   │   │   ├── memory_cache.go       # In-memory cache
│   │   │   ├── file_cache.go         # File-based cache
│   │   │   ├── interfaces.go         # Cache interfaces
│   │   │   └── ttl_cache.go          # TTL cache
│   │   │
│   │   ├── utils/                    # Utility Functions
│   │   │   ├── strings.go            # String utilities
│   │   │   ├── time.go               # Time utilities
│   │   │   ├── path.go               # Path utilities
│   │   │   └── version.go            # Version utilities
│   │   │
│   │   └── patterns/                 # Design Patterns
│   │       ├── observer.go           # Observer pattern
│   │       ├── factory.go            # Factory pattern
│   │       └── singleton.go          # Singleton pattern
│   │
│   └── tests/                        # Test Utilities and Fixtures
│       ├── mocks/                    # Generated and manual mocks
│       │   ├── generated/            # Auto-generated mocks (gomock)
│       │   └── manual/               # Manually created mocks
│       ├── fixtures/                 # Test data and fixtures
│       │   ├── configs/              # Test configurations
│       │   ├── labs/                 # Test lab files
│       │   └── manifests/            # Test manifests
│       ├── helpers/                  # Test helper functions
│       │   ├── assertions.go         # Custom assertions
│       │   ├── setup.go              # Test setup utilities
│       │   └── teardown.go           # Test cleanup
│       └── integration/              # Integration test utilities
│           ├── docker_helper.go      # Docker test helpers
│           ├── k8s_helper.go         # Kubernetes test helpers
│           └── cluster_helper.go     # Cluster test helpers
│
├── configs/                          # Configuration Files
│   ├── default.yaml                  # Default configuration
│   ├── development.yaml              # Development config
│   ├── production.yaml               # Production config
│   └── examples/                     # Example configurations
│       ├── minimal.yaml
│       └── advanced.yaml
│
├── deployments/                      # Deployment Manifests and Scripts
│   ├── docker/                       # Docker deployments
│   │   ├── Dockerfile
│   │   ├── docker-compose.yaml
│   │   └── .dockerignore
│   ├── k8s/                         # Kubernetes manifests
│   │   ├── namespace.yaml
│   │   ├── deployment.yaml
│   │   ├── service.yaml
│   │   └── rbac.yaml
│   └── scripts/                     # Deployment scripts
│       ├── deploy.sh
│       └── rollback.sh
│
├── docs/                            # Documentation
│   ├── architecture/                # Architecture documentation
│   │   ├── README.md
│   │   ├── adr/                     # Architecture Decision Records
│   │   └── diagrams/                # Architecture diagrams
│   ├── api/                         # API documentation
│   ├── user/                        # User documentation
│   │   ├── installation.md
│   │   ├── quickstart.md
│   │   └── tutorials/
│   └── development/                 # Development documentation
│       ├── contributing.md
│       ├── testing.md
│       └── release.md
│
├── scripts/                         # Build and Development Scripts
│   ├── build.sh                     # Build script
│   ├── test.sh                      # Test runner
│   ├── lint.sh                      # Linting script
│   ├── release.sh                   # Release script
│   ├── install.sh                   # Installation script
│   └── dev/                         # Development scripts
│       ├── setup.sh                 # Development setup
│       ├── generate-mocks.sh        # Mock generation
│       └── update-deps.sh           # Dependency updates
│
├── tests/                           # End-to-End and Integration Tests
│   ├── e2e/                         # End-to-end tests
│   │   ├── cluster_test.go
│   │   ├── lab_test.go
│   │   └── fixtures/
│   ├── integration/                 # Integration tests
│   │   ├── repository_test.go
│   │   ├── container_engine_test.go
│   │   └── k8s_test.go
│   └── performance/                 # Performance tests
│       ├── cluster_creation_test.go
│       └── benchmarks/
│
├── tools/                           # Development Tools
│   ├── mockgen/                     # Mock generation tools
│   │   └── generate.go
│   ├── linter/                      # Custom linters
│   └── migrate/                     # Migration tools
│
├── vendor/                          # Vendored dependencies (if using)
│
├── .github/                         # GitHub specific files
│   ├── workflows/                   # GitHub Actions
│   │   ├── ci.yml
│   │   ├── release.yml
│   │   └── security.yml
│   ├── ISSUE_TEMPLATE/              # Issue templates
│   └── PULL_REQUEST_TEMPLATE.md     # PR template
│
├── .gitignore                       # Git ignore rules
├── .golangci.yml                    # Linter configuration
├── go.mod                           # Go modules
├── go.sum                           # Go modules checksum
├── Makefile                         # Build automation
├── LICENSE                          # License file
└── README.md                        # Project documentation

```

## 📁 Detalhamento das Camadas

### 7.1 Application Layer (`internal/app/`)

**Responsabilidade**: Orquestração de casos de uso e coordenação entre domínios.

```txt
internal/app/handlers/cluster/create_handler.go
```

```go
type CreateHandler struct {
    orchestrator *orchestrators.ClusterOrchestrator
    validator    *validation.Validator
    logger       logger.Logger
}

func (h *CreateHandler) Handle(ctx context.Context, req *dto.CreateClusterRequest) (*dto.CreateClusterResponse, error) {
    // 1. Validar request
    if err := h.validator.ValidateStruct(req); err != nil {
        return nil, errors.NewValidationError("invalid request", err)
    }
    
    // 2. Converter para domain request
    domainReq := h.toDomainRequest(req)
    
    // 3. Executar através do orchestrator
    result, err := h.orchestrator.CreateCluster(ctx, domainReq)
    if err != nil {
        return nil, err
    }
    
    // 4. Converter response
    return h.toResponse(result), nil
}
```

### 7.2 Domain Layer (`internal/domain/`)

**Responsabilidade**: Lógica de negócio pura e regras de domínio.

```txt
internal/domain/cluster/entities.go
```

```go
type Cluster struct {
    id          ClusterID
    name        string
    status      ClusterStatus
    nodes       []*Node
    createdAt   time.Time
    updatedAt   time.Time
}

func NewCluster(name string, config *CreationConfig) (*Cluster, error) {
    if err := validateClusterName(name); err != nil {
        return nil, err
    }
    
    return &Cluster{
        id:        NewClusterID(),
        name:      name,
        status:    StatusCreating,
        nodes:     make([]*Node, 0),
        createdAt: time.Now(),
        updatedAt: time.Now(),
    }, nil
}

func (c *Cluster) AddNode(node *Node) error {
    if c.status != StatusReady {
        return errors.New("cannot add node to non-ready cluster")
    }
    
    c.nodes = append(c.nodes, node)
    c.updatedAt = time.Now()
    return nil
}
```

### 7.3 Infrastructure Layer (`internal/infrastructure/`)

**Responsabilidade**: Implementações técnicas e detalhes de infraestrutura.

```txt
internal/infrastructure/k8s/cluster_repository.go
```

```go
type ClusterRepository struct {
    client    kubernetes.Interface
    logger    logger.Logger
    namespace string
}

func (r *ClusterRepository) Save(ctx context.Context, cluster *domain.Cluster) error {
    configMap := &corev1.ConfigMap{
        ObjectMeta: metav1.ObjectMeta{
            Name:      fmt.Sprintf("cluster-%s", cluster.Name()),
            Namespace: r.namespace,
            Labels: map[string]string{
                "app.kubernetes.io/name":      "girus",
                "app.kubernetes.io/component": "cluster",
            },
        },
        Data: map[string]string{
            "cluster.json": r.serializeCluster(cluster),
        },
    }
    
    _, err := r.client.CoreV1().ConfigMaps(r.namespace).Create(ctx, configMap, metav1.CreateOptions{})
    return err
}
```

### 7.4 Interface Adapters (`internal/adapters/`)

**Responsabilidade**: Conversão entre camadas e adaptação de interfaces externas.

```txt
internal/adapters/config/loader.go
```

```go
type ConfigLoader struct {
    sources []ConfigSource
    merger  *ConfigMerger
}

func (l *ConfigLoader) Load() (*Config, error) {
    var configs []*Config
    
    for _, source := range l.sources {
        if config, err := source.Load(); err == nil {
            configs = append(configs, config)
        }
    }
    
    return l.merger.Merge(configs...)
}
```

## 🔄 Migração Gradual

### 7.5 Estratégia de Migração

1. **Fase 1**: Criar estrutura de pastas
2. **Fase 2**: Extrair interfaces de domínio
3. **Fase 3**: Implementar services de domínio
4. **Fase 4**: Criar handlers e orchestrators
5. **Fase 5**: Migrar implementações de infraestrutura
6. **Fase 6**: Adicionar adapters e configuração
7. **Fase 7**: Refatorar comandos CLI

### 7.6 Exemplo de Migração: Create Cluster

**Antes** (arquivo atual de 500+ linhas):

```txt
cmd/create.go
```

```go
var createClusterCmd = &cobra.Command{
    Use: "cluster",
    Run: func(cmd *cobra.Command, args []string) {
        // 500+ linhas de lógica misturada
    },
}

```

**Depois** (estrutura limpa):

```txt
cmd/create.go
```

```go
var createClusterCmd = &cobra.Command{
    Use: "cluster",
    RunE: func(cmd *cobra.Command, args []string) error {
        req := buildCreateClusterRequest(cmd)
        handler := app.Container.GetCreateClusterHandler()
        return handler.Handle(cmd.Context(), req)
    },
}

// internal/app/handlers/cluster/create_handler.go - 50 linhas
// internal/app/orchestrators/cluster_orchestrator.go - 80 linhas
// internal/domain/cluster/service.go - 60 linhas
// internal/infrastructure/kind/cluster_manager.go - 70 linhas

```

## 📊 Benefícios da Nova Estrutura

| Aspecto | Antes | Depois |
|---------|-------|---------|
| **Organização** | Funcional (por tipo) | Por domínio e responsabilidade |
| **Acoplamento** | Alto (tudo misturado) | Baixo (camadas bem definidas) |
| **Testabilidade** | Difícil (dependências hardcoded) | Fácil (interfaces mockáveis) |
| **Reutilização** | Baixa (código duplicado) | Alta (componentes modulares) |
| **Manutenção** | Complexa (mudanças em cascata) | Simples (mudanças isoladas) |
| **Escalabilidade** | Limitada | Excelente (fácil adicionar domínios) |

## 🎯 Princípios Seguidos

1. **Single Responsibility Principle**: Cada pasta/arquivo tem uma responsabilidade
2. **Dependency Inversion**: Dependências apontam para abstrações
3. **Interface Segregation**: Interfaces específicas e focadas
4. **Domain-Driven Design**: Organização por domínios de negócio
5. **Clean Architecture**: Separação clara de camadas
6. **Convention over Configuration**: Estrutura padronizada e previsível

## 🔄 Próxima Etapa

[Fluxograma do Novo Funcionamento](./08-fluxograma.md) - Visualização do fluxo de execução.
