# 4. Interfaces e Contratos

## 🎯 Objetivo

Definir interfaces claras e bem estruturadas que permitam inversão de dependência, facilitando testes e garantindo flexibilidade para diferentes implementações.

## 🏗️ Princípios de Design

### 4.1 Interface Segregation Principle (ISP)

Interfaces específicas e focadas, evitando dependências desnecessárias:

```go

// ❌ Interface muito ampla (viola ISP)
type ClusterService interface {
    Create(ctx context.Context, config *CreateConfig) error
    Delete(ctx context.Context, name string) error
    List(ctx context.Context) ([]*Cluster, error)
    Deploy(ctx context.Context, manifests []string) error
    GetLogs(ctx context.Context, pod string) ([]string, error)
    PortForward(ctx context.Context, service string, port int) error
}

// ✅ Interfaces segregadas (seguem ISP)
type ClusterManager interface {
    Create(ctx context.Context, config *CreateConfig) error
    Delete(ctx context.Context, name string) error
    Exists(ctx context.Context, name string) (bool, error)
    GetStatus(ctx context.Context, name string) (*ClusterStatus, error)
}

type ClusterDeployer interface {
    Deploy(ctx context.Context, config *DeployConfig) error
    Validate(manifests []string) (*ValidationResult, error)
}

type ClusterAccess interface {
    SetupPortForward(ctx context.Context, config *PortForwardConfig) error
    OpenBrowser(ctx context.Context, url string) error
}

```

### 4.2 Dependency Inversion Principle (DIP)

Depender de abstrações, não de concretizações:

```txt
internal/domain/cluster/interfaces.go
```

```go
type Repository interface {
    Save(ctx context.Context, cluster *Cluster) error
    FindByName(ctx context.Context, name string) (*Cluster, error)
    List(ctx context.Context, filters *ListFilters) ([]*Cluster, error)
    Delete(ctx context.Context, name string) error
}

type ContainerEngine interface {
    IsAvailable(ctx context.Context) error
    IsRunning(ctx context.Context) error
    GetVersion(ctx context.Context) (*Version, error)
    Start(ctx context.Context) error
}

type KubernetesClient interface {
    ApplyManifests(ctx context.Context, manifests []string) error
    WaitForPods(ctx context.Context, namespace string, timeout time.Duration) error
    GetPodStatus(ctx context.Context, namespace, name string) (*PodStatus, error)
}
```

## 📋 Interfaces por Domínio

### 4.3 Cluster Domain

```txt
internal/domain/cluster/service.go
```

```go
type Service interface {
    Create(ctx context.Context, config *CreateConfig) (*Cluster, error)
    Delete(ctx context.Context, name string) error
    GetByName(ctx context.Context, name string) (*Cluster, error)
    List(ctx context.Context, filters *ListFilters) ([]*Cluster, error)
}

type CreateConfig struct {
    Name            string                 `validate:"required,min=3,max=63"`
    ContainerEngine string                 `validate:"required,oneof=docker podman"`
    KubeConfig      string                 `validate:"omitempty,file"`
    Resources       *ResourceRequirements  `validate:"omitempty"`
    Networking      *NetworkConfig         `validate:"omitempty"`
    Labels          map[string]string      `validate:"omitempty"`
    Annotations     map[string]string      `validate:"omitempty"`
}

type Cluster struct {
    Name        string
    Status      ClusterStatus
    CreatedAt   time.Time
    UpdatedAt   time.Time
    Nodes       []*Node
    Services    []*Service
    Resources   *AllocatedResources
    Metadata    *ClusterMetadata
}

type ClusterStatus string

const (
    ClusterStatusCreating ClusterStatus = "creating"
    ClusterStatusReady    ClusterStatus = "ready"
    ClusterStatusFailed   ClusterStatus = "failed"
    ClusterStatusDeleting ClusterStatus = "deleting"
)
```

### 4.4 Lab Domain

```txt
internal/domain/lab/service.go
```

```go
type Service interface {
    Install(ctx context.Context, req *InstallRequest) (*Lab, error)
    Uninstall(ctx context.Context, labID string) error
    Get(ctx context.Context, labID string) (*Lab, error)
    List(ctx context.Context, filters *ListFilters) ([]*Lab, error)
    Validate(content []byte) (*ValidationResult, error)
}

type InstallRequest struct {
    Source      LabSource              `validate:"required"`
    Content     []byte                 `validate:"required_if=Source file"`
    RepositoryURL string               `validate:"required_if=Source repo,omitempty,url"`
    LabID       string                 `validate:"required_if=Source repo"`
    Version     string                 `validate:"omitempty,semver"`
    Options     *InstallOptions        `validate:"omitempty"`
}

type Lab struct {
    ID          string
    Name        string
    Version     string
    Description string
    Author      string
    Tags        []string
    Status      LabStatus
    Spec        *LabSpec
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

type LabSpec struct {
    Environment *EnvironmentSpec `yaml:"environment"`
    Tasks       []*Task          `yaml:"tasks"`
    Validation  *ValidationSpec  `yaml:"validation"`
    Resources   *ResourceSpec    `yaml:"resources"`
}

type Validator interface {
    Validate(lab *Lab) (*ValidationResult, error)
    ValidateEnvironment(env *EnvironmentSpec) []ValidationError
    ValidateTasks(tasks []*Task) []ValidationError
    ValidateResources(resources *ResourceSpec) []ValidationError
}
```

### 4.5 Infrastructure Domain

```txt
internal/domain/infrastructure/service.go
```

```go
type Service interface {
    CheckPrerequisites(ctx context.Context, req *PrerequisiteCheck) (*CheckResult, error)
    InstallMissingTools(ctx context.Context, tools []string) error
    GetSystemInfo(ctx context.Context) (*SystemInfo, error)
}

type PrerequisiteChecker interface {
    CheckContainerEngine(ctx context.Context, engine string) error
    CheckKubernetes(ctx context.Context) error
    CheckPorts(ctx context.Context, ports []int) error
    CheckResources(ctx context.Context, requirements *ResourceRequirements) error
}

type ToolInstaller interface {
    Install(ctx context.Context, tool *Tool) error
    IsInstalled(tool string) bool
    GetVersion(tool string) (*Version, error)
}

type SystemDetector interface {
    GetOS() string
    GetArch() string
    GetDistribution() string
    GetKernelVersion() string
    GetAvailableResources() *AvailableResources
}
```

### 4.6 UI/Presentation Domain

```txt
internal/domain/ui/service.go
```

```go
type Service interface {
    ShowHeader(title string)
    ShowMessage(level MessageLevel, message string, details ...string)
    ShowProgress(ctx context.Context, config *ProgressConfig) (ProgressBar, error)
    Confirm(message string, defaultValue bool) (bool, error)
    SelectOption(prompt string, options []SelectOption) (string, error)
    ShowTable(headers []string, rows [][]string)
}

type MessageLevel string

const (
    MessageLevelInfo    MessageLevel = "info"
    MessageLevelSuccess MessageLevel = "success"
    MessageLevelWarning MessageLevel = "warning"
    MessageLevelError   MessageLevel = "error"
)

type ProgressBar interface {
    Start() error
    Update(current int, description string) error
    Finish() error
    Failed(err error) error
}

type Formatter interface {
    FormatError(err error) string
    FormatSuccess(message string) string
    FormatWarning(message string) string
    FormatHeader(title string) string
}

type ColorProvider interface {
    Red(text string) string
    Green(text string) string
    Yellow(text string) string
    Blue(text string) string
    Magenta(text string) string
    Cyan(text string) string
    Bold(text string) string
}
```

## 🔧 Implementações Concretas

### 4.7 Exemplo de Implementação

```txt
internal/services/cluster/kind_service.go
```

```go
type KindClusterService struct {
    containerEngine ContainerEngine
    k8sClient      KubernetesClient
    repository     Repository
    validator      Validator
    logger         logger.Logger
}

// Garantir que implementa a interface
var _ cluster.Service = (*KindClusterService)(nil)

func NewKindClusterService(
    containerEngine ContainerEngine,
    k8sClient KubernetesClient,
    repository Repository,
    validator Validator,
    logger logger.Logger,
) *KindClusterService {
    return &KindClusterService{
        containerEngine: containerEngine,
        k8sClient:      k8sClient,
        repository:     repository,
        validator:      validator,
        logger:         logger,
    }
}

func (s *KindClusterService) Create(ctx context.Context, config *CreateConfig) (*Cluster, error) {
    // 1. Validar configuração
    if err := s.validator.ValidateCreateConfig(config); err != nil {
        return nil, fmt.Errorf("invalid config: %w", err)
    }
    
    // 2. Verificar se cluster já existe
    existing, _ := s.repository.FindByName(ctx, config.Name)
    if existing != nil {
        return nil, domain.NewClusterExistsError(config.Name)
    }
    
    // 3. Verificar container engine
    if err := s.containerEngine.IsAvailable(ctx); err != nil {
        return nil, fmt.Errorf("container engine not available: %w", err)
    }
    
    // 4. Criar cluster
    cluster, err := s.createKindCluster(ctx, config)
    if err != nil {
        return nil, fmt.Errorf("failed to create cluster: %w", err)
    }
    
    // 5. Salvar no repositório
    if err := s.repository.Save(ctx, cluster); err != nil {
        s.logger.Warn("failed to save cluster to repository", "error", err)
    }
    
    return cluster, nil
}
```

## 🧪 Mock Interfaces para Testes

### 4.8 Geração Automática de Mocks

```bash
go:generate mockgen -source=service.go -destination=mocks/service_mock.go

internal/domain/cluster/mocks/service_mock.go (gerado automaticamente)
```

```go
type MockService struct {
    ctrl     *gomock.Controller
    recorder *MockServiceMockRecorder
}

func (m *MockService) Create(ctx context.Context, config *CreateConfig) (*Cluster, error) {
    m.ctrl.T.Helper()
    ret := m.ctrl.Call(m, "Create", ctx, config)
    ret0, _ := ret[0].(*Cluster)
    ret1, _ := ret[1].(error)
    return ret0, ret1
}

// Uso em testes
func TestClusterOrchestrator_CreateCluster(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    
    mockClusterSvc := mocks.NewMockService(ctrl)
    mockPrereqSvc := mocks.NewMockPrerequisiteService(ctrl)
    
    orchestrator := services.NewClusterOrchestrator(
        mockClusterSvc,
        mockPrereqSvc,
        /* ... */
    )
    
    // Setup expectations
    mockPrereqSvc.EXPECT().
        CheckPrerequisites(gomock.Any(), gomock.Any()).
        Return(&CheckResult{AllPassed: true}, nil)
    
    mockClusterSvc.EXPECT().
        Create(gomock.Any(), gomock.Any()).
        Return(&Cluster{Name: "test"}, nil)
    
    // Execute test
    result, err := orchestrator.CreateCluster(context.Background(), &CreateClusterRequest{
        Name: "test-cluster",
    })
    
    assert.NoError(t, err)
    assert.Equal(t, "test", result.Name)
}
```

## 📏 Validação de Contratos

### 4.9 Interface Compliance Tests

```txt
internal/tests/interfaces/compliance_test.go
```

```go
func TestInterfaceCompliance(t *testing.T) {
    // Verificar se todas as implementações seguem os contratos
    var _ cluster.Service = (*cluster_kind.Service)(nil)
    var _ cluster.Service = (*cluster_k3s.Service)(nil)
    
    var _ lab.Service = (*lab_configmap.Service)(nil)
    var _ lab.Service = (*lab_helm.Service)(nil)
    
    var _ infrastructure.Service = (*infrastructure_linux.Service)(nil)
    var _ infrastructure.Service = (*infrastructure_darwin.Service)(nil)
}

func TestContractBehavior(t *testing.T) {
    // Testes que verificam se implementações seguem contratos comportamentais
    implementations := []cluster.Service{
        cluster_kind.NewService(/* deps */),
        cluster_k3s.NewService(/* deps */),
    }
    
    for _, impl := range implementations {
        t.Run(fmt.Sprintf("%T", impl), func(t *testing.T) {
            // Teste de comportamento padrão
            _, err := impl.Create(context.Background(), &cluster.CreateConfig{
                Name: "test-cluster",
            })
            
            // Todos devem retornar erro para nome vazio
            assert.Error(t, err)
        })
    }
}
```

## 📊 Benefícios das Interfaces Bem Definidas

| Aspecto | Benefício | Exemplo |
|---------|-----------|---------|
| **Testabilidade** | Mocks facilmente criados | `MockClusterService` |
| **Flexibilidade** | Múltiplas implementações | Kind, K3s, K8s |
| **Manutenibilidade** | Mudanças isoladas | Trocar Docker por Podman |
| **Documentação** | Contrato claro | Interface documenta comportamento |
| **Evolução** | Mudanças compatíveis | Adicionar métodos opcionais |

## 🔄 Próxima Etapa

[Sistema de Configuração](./05-sistema-configuracao.md) - Configuração robusta e flexível.
