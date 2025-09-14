# 2. Arquitetura de Services

## 🏗️ Visão Geral

A nova arquitetura propõe a criação de uma camada de services que separa a lógica de negócio dos comandos CLI, seguindo princípios de Clean Architecture e inversão de dependência.

## 📋 Services Propostos

### 2.1 ClusterService

**Responsabilidade**: Gerenciar operações relacionadas a clusters Kubernetes.

```txt
internal/domain/cluster/service.go
```

```go
type Service interface {
    Create(ctx context.Context, config *CreateConfig) error
    Delete(ctx context.Context, name string) error
    Exists(ctx context.Context, name string) (bool, error)
    GetStatus(ctx context.Context, name string) (*Status, error)
    Deploy(ctx context.Context, config *DeployConfig) error
    SetupAccess(ctx context.Context, config *AccessConfig) error
}

type CreateConfig struct {
    Name            string
    ContainerEngine string
    KubeConfig      string
    DeployFiles     []string
    Resources       *ResourceConfig
    Networking      *NetworkConfig
    Verbose         bool
}

type DeployConfig struct {
    ClusterName string
    Manifests   []string
    Templates   []TemplateConfig
    WaitReady   bool
    Timeout     time.Duration
}

type AccessConfig struct {
    ClusterName     string
    PortForward     bool
    OpenBrowser     bool
    Services        []ServicePort
}
```

**Implementação Sugerida**:

```txt
internal/services/cluster/kind_service.go
```

```go
type KindService struct {
    executor      exec.CommandExecutor
    k8sClient     k8s.Client
    validator     config.Validator
    logger        logger.Logger
    progressSvc   progress.Service
}

func (s *KindService) Create(ctx context.Context, config *CreateConfig) error {
    // 1. Validar configuração
    if err := s.validator.ValidateCreateConfig(config); err != nil {
        return fmt.Errorf("invalid config: %w", err)
    }
    
    // 2. Verificar se cluster já existe
    exists, err := s.Exists(ctx, config.Name)
    if err != nil {
        return fmt.Errorf("failed to check cluster existence: %w", err)
    }
    
    if exists {
        return domain.NewClusterExistsError(config.Name)
    }
    
    // 3. Criar cluster
    return s.createKindCluster(ctx, config)
}

func (s *KindService) createKindCluster(ctx context.Context, config *CreateConfig) error {
    progress := s.progressSvc.NewProgressBar(&progress.Config{
        Title: "Creating cluster",
        Total: 100,
    })
    defer progress.Finish()
    
    // Comando kind create cluster
    cmd := s.executor.Command("kind", "create", "cluster", "--name", config.Name)
    
    return s.executor.RunWithProgress(ctx, cmd, progress)
}
```

### 2.2 LabService

**Responsabilidade**: Gerenciar laboratórios e templates.

```txt
internal/domain/lab/service.go
```

```go
type Service interface {
    Install(ctx context.Context, req *InstallRequest) (*Lab, error)
    Uninstall(ctx context.Context, labID string) error
    List(ctx context.Context, filters *ListFilters) ([]*Lab, error)
    Get(ctx context.Context, labID string) (*Lab, error)
    Validate(labFile string) (*ValidationResult, error)
}

type InstallRequest struct {
    Source   LabSource
    FilePath string
    RepoURL  string
    LabID    string
    Options  *InstallOptions
}

type Lab struct {
    ID          string
    Name        string
    Version     string
    Description string
    Status      LabStatus
    Metadata    *LabMetadata
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

type ValidationResult struct {
    Valid    bool
    Errors   []ValidationError
    Warnings []ValidationWarning
}
```

### 2.3 PrerequisiteService

**Responsabilidade**: Validar pré-requisitos do sistema.

```txt
internal/domain/infrastructure/service.go
```

```go
type Service interface {
    CheckAll(ctx context.Context, config *CheckConfig) error
    CheckContainerEngine(ctx context.Context, engine string) error
    CheckKubernetes(ctx context.Context) error
    GetSystemInfo(ctx context.Context) (*SystemInfo, error)
    InstallMissing(ctx context.Context, tools []string) error
}

type CheckConfig struct {
    ContainerEngine string
    RequiredTools   []string
    RequiredPorts   []int
    Resources       *ResourceRequirements
}

type SystemInfo struct {
    OS              string
    Arch            string
    ContainerEngine *ContainerEngineInfo
    Kubernetes      *KubernetesInfo
    AvailablePorts  []int
    Resources       *AvailableResources
}
```

**Implementação**:

```txt
internal/services/infrastructure/prerequisite_service.go
```

```go
type PrerequisiteService struct {
    executor exec.CommandExecutor
    detector system.Detector
    logger   logger.Logger
}

func (s *PrerequisiteService) CheckContainerEngine(ctx context.Context, engine string) error {
    // 1. Verificar se está instalado
    if !s.detector.IsInstalled(engine) {
        return NewContainerEngineNotFoundError(engine)
    }
    
    // 2. Verificar se está rodando
    if !s.detector.IsRunning(engine) {
        return NewContainerEngineNotRunningError(engine)
    }
    
    // 3. Testar conectividade
    return s.testContainerEngineConnectivity(ctx, engine)
}
```

### 2.4 ProgressService

**Responsabilidade**: Gerenciar interfaces de progresso.

```txt
internal/domain/ui/service.go
```

```go
type Service interface {
    ShowHeader(title string)
    ShowProgress(ctx context.Context, config *ProgressConfig) ProgressBar
    ShowSuccess(message string, details ...string)
    ShowError(err error) error
    ShowWarning(message string, details ...string)
    Confirm(message string) (bool, error)
    SelectOption(prompt string, options []string) (string, error)
}

type ProgressConfig struct {
    Title       string
    Description string
    Total       int
    ShowBytes   bool
    ShowTime    bool
    Style       ProgressStyle
}

type ProgressBar interface {
    Add(int) error
    Set(int) error
    Finish() error
    Describe(string)
}
```

## 🔗 Integração entre Services

### Exemplo: Create Cluster Flow

```txt
internal/app/services/cluster_orchestrator.go
```

```go
type ClusterOrchestrator struct {
    clusterSvc      cluster.Service
    prerequisiteSvc infrastructure.Service
    labSvc          lab.Service
    uiSvc           ui.Service
    config          *config.Config
}

func (o *ClusterOrchestrator) CreateCluster(ctx context.Context, req *CreateClusterRequest) error {
    // 1. Mostrar cabeçalho
    o.uiSvc.ShowHeader("GIRUS CREATE")
    
    // 2. Verificar pré-requisitos
    if err := o.prerequisiteSvc.CheckAll(ctx, &infrastructure.CheckConfig{
        ContainerEngine: req.ContainerEngine,
        RequiredTools:   []string{"kind", "kubectl"},
    }); err != nil {
        return o.handlePrerequisiteError(err)
    }
    
    // 3. Criar cluster
    clusterConfig := &cluster.CreateConfig{
        Name:            req.Name,
        ContainerEngine: req.ContainerEngine,
        Resources:       req.Resources,
    }
    
    if err := o.clusterSvc.Create(ctx, clusterConfig); err != nil {
        return fmt.Errorf("failed to create cluster: %w", err)
    }
    
    // 4. Deploy infraestrutura
    deployConfig := &cluster.DeployConfig{
        ClusterName: req.Name,
        Manifests:   req.DeployFiles,
        WaitReady:   true,
        Timeout:     5 * time.Minute,
    }
    
    if err := o.clusterSvc.Deploy(ctx, deployConfig); err != nil {
        return fmt.Errorf("failed to deploy infrastructure: %w", err)
    }
    
    // 5. Configurar acesso
    if !req.SkipAccess {
        accessConfig := &cluster.AccessConfig{
            ClusterName: req.Name,
            PortForward: !req.SkipPortForward,
            OpenBrowser: !req.SkipBrowser,
        }
        
        return o.clusterSvc.SetupAccess(ctx, accessConfig)
    }
    
    return nil
}
```

## 🧪 Benefícios da Arquitetura de Services

### Testabilidade

```go
func TestClusterOrchestrator_CreateCluster(t *testing.T) {
    mockClusterSvc := &cluster.MockService{}
    mockPrereqSvc := &infrastructure.MockService{}
    mockUISvc := &ui.MockService{}
    
    orchestrator := &ClusterOrchestrator{
        clusterSvc:      mockClusterSvc,
        prerequisiteSvc: mockPrereqSvc,
        uiSvc:           mockUISvc,
    }
    
    // Setup mocks
    mockPrereqSvc.On("CheckAll", mock.Anything, mock.Anything).Return(nil)
    mockClusterSvc.On("Create", mock.Anything, mock.Anything).Return(nil)
    
    // Test
    err := orchestrator.CreateCluster(context.Background(), &CreateClusterRequest{
        Name: "test-cluster",
    })
    
    assert.NoError(t, err)
    mockClusterSvc.AssertExpectations(t)
}
```

### Reutilização

- Services podem ser usados por diferentes interfaces (CLI, Web, API)
- Componentes modulares e configuráveis
- Lógica de negócio independente de apresentação

### Manutenibilidade

- Responsabilidades claras e bem definidas
- Fácil localização de bugs
- Mudanças isoladas e controladas

## 📈 Comparação: Antes vs Depois

| Aspecto | Antes | Depois |
|---------|-------|---------|
| Linhas de código por função | 500+ | < 50 |
| Responsabilidades | Múltiplas | Única |
| Testabilidade | Impossível | Excelente |
| Reutilização | Baixa | Alta |
| Manutenibilidade | Difícil | Fácil |

## 🔄 Próxima Etapa

[Padrão Command/Handler](./03-padrao-command-handler.md) - Como os comandos CLI interagem com os services.
