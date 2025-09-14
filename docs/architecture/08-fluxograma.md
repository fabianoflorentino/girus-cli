# 8. Fluxograma do Novo Funcionamento

## 🎯 Objetivo

Visualizar o fluxo de execução da nova arquitetura através de diagramas que mostram como os diferentes componentes interagem para executar operações como criação de cluster e instalação de laboratórios.

## 🔄 Fluxo Principal: Create Cluster

### 8.1 Diagrama de Sequência - Create Cluster

```mermaid
sequenceDiagram
    participant User as 👤 User
    participant CLI as 📱 CLI Command
    participant Handler as 🎯 Create Handler
    participant Orch as 🎼 Cluster Orchestrator
    participant PreReq as ⚡ Prerequisites Service
    participant Cluster as 🏗️ Cluster Service
    participant K8s as ☸️ K8s Service
    participant UI as 🎨 UI Service
    participant Config as ⚙️ Config Service

    User->>CLI: girus create cluster --name mylab
    
    CLI->>Config: Load configuration
    Config-->>CLI: Merged config (file + env + flags)
    
    CLI->>Handler: Handle(ctx, CreateClusterRequest)
    
    Handler->>Handler: Validate request
    
    Handler->>Orch: CreateCluster(ctx, request)
    
    Orch->>UI: ShowHeader("GIRUS CREATE")
    UI-->>User: 🎨 Display header
    
    Orch->>PreReq: CheckPrerequisites(containerEngine)
    PreReq->>PreReq: Check Docker/Podman
    PreReq->>PreReq: Check Kind/kubectl
    PreReq->>PreReq: Check available resources
    
    alt Prerequisites Failed
        PreReq-->>Orch: PrerequisiteError
        Orch->>UI: ShowError + Suggestions
        UI-->>User: ❌ Error with recovery steps
        Orch-->>Handler: Error
        Handler-->>CLI: Error
        CLI-->>User: Exit with error code
    else Prerequisites OK
        PreReq-->>Orch: Success ✅
        
        Orch->>Cluster: Create(ctx, ClusterConfig)
        
        Cluster->>Cluster: Validate cluster name
        Cluster->>Cluster: Check if cluster exists
        
        alt Cluster Exists
            Cluster->>UI: Confirm("Replace existing cluster?")
            UI-->>User: ❓ Confirmation prompt
            User-->>UI: User response
            UI-->>Cluster: Boolean response
            
            alt User Confirms
                Cluster->>K8s: DeleteCluster(name)
                K8s->>UI: ShowProgress("Deleting cluster...")
                UI-->>User: 📊 Progress bar
                K8s-->>Cluster: Success
            else User Denies
                Cluster-->>Orch: OperationCancelledError
                Orch-->>Handler: Error
                Handler-->>CLI: Error
                CLI-->>User: Operation cancelled
            end
        end
        
        Cluster->>K8s: CreateKindCluster(config)
        K8s->>UI: ShowProgress("Creating cluster...")
        UI-->>User: 📊 Progress bar
        K8s->>K8s: Execute kind create cluster
        K8s-->>Cluster: Cluster created ✅
        
        Cluster->>K8s: DeployInfrastructure(manifests)
        K8s->>UI: ShowProgress("Deploying infrastructure...")
        UI-->>User: 📊 Progress bar
        K8s->>K8s: Apply YAML manifests
        K8s->>K8s: Wait for pods ready
        K8s-->>Cluster: Infrastructure deployed ✅
        
        Cluster->>K8s: SetupAccess(portForward, browser)
        K8s->>K8s: Setup port-forwarding
        K8s->>K8s: Open browser (if enabled)
        K8s-->>Cluster: Access configured ✅
        
        Cluster-->>Orch: ClusterCreateResult
        Orch->>UI: ShowSuccess("Cluster ready!")
        UI-->>User: ✅ Success message
        
        Orch-->>Handler: Success
        Handler-->>CLI: Success
        CLI-->>User: Exit code 0
    end

```

### 8.2 Diagrama de Fluxo - Create Cluster

```mermaid

flowchart TD
    Start([👤 User runs: girus create cluster]) --> LoadConfig[⚙️ Load Config<br/>File + Env + Flags]
    
    LoadConfig --> ParseFlags[📝 Parse CLI Flags<br/>--name, --engine, --verbose]
    
    ParseFlags --> ValidateInput{✅ Validate Input<br/>Required fields, formats}
    
    ValidateInput -->|Invalid| ShowInputError[❌ Show Input Error<br/>+ Suggestions]
    ShowInputError --> Exit1([🚪 Exit 1])
    
    ValidateInput -->|Valid| CreateHandler[🎯 Create Handler<br/>Handle request]
    
    CreateHandler --> ShowHeader[🎨 Show Header<br/>GIRUS CREATE]
    
    ShowHeader --> CheckPrereqs[⚡ Check Prerequisites]
    
    CheckPrereqs --> CheckContainer{🐳 Container Engine<br/>Available?}
    
    CheckContainer -->|No| TryRecover{🔧 Auto Recovery<br/>Possible?}
    TryRecover -->|Yes| RecoverEngine[🔧 Start/Install Engine]
    RecoverEngine --> CheckContainer
    TryRecover -->|No| ShowEngineError[❌ Show Engine Error<br/>+ Install instructions]
    ShowEngineError --> Exit2([🚪 Exit 1])
    
    CheckContainer -->|Yes| CheckK8sTools{☸️ K8s Tools<br/>Available?}
    
    CheckK8sTools -->|No| ShowK8sError[❌ Show K8s Tools Error<br/>+ Install instructions]
    ShowK8sError --> Exit3([🚪 Exit 1])
    
    CheckK8sTools -->|Yes| CheckResources{💾 Resources<br/>Available?}
    
    CheckResources -->|No| ShowResourceError[❌ Show Resource Error<br/>+ Suggestions]
    ShowResourceError --> Exit4([🚪 Exit 1])
    
    CheckResources -->|Yes| CheckClusterExists{🔍 Cluster<br/>Exists?}
    
    CheckClusterExists -->|Yes| ConfirmReplace{❓ Replace<br/>Existing?}
    ConfirmReplace -->|No| ShowCancelError[⚠️ Operation Cancelled]
    ShowCancelError --> Exit5([🚪 Exit 0])
    
    ConfirmReplace -->|Yes| DeleteExisting[🗑️ Delete Existing<br/>+ Progress Bar]
    DeleteExisting --> CreateCluster
    
    CheckClusterExists -->|No| CreateCluster[🏗️ Create Cluster<br/>+ Progress Bar]
    
    CreateCluster --> ClusterReady{✅ Cluster<br/>Created?}
    
    ClusterReady -->|No| ShowCreateError[❌ Show Create Error<br/>+ Debug info]
    ShowCreateError --> Exit6([🚪 Exit 1])
    
    ClusterReady -->|Yes| DeployInfra[📦 Deploy Infrastructure<br/>+ Progress Bar]
    
    DeployInfra --> InfraReady{✅ Infrastructure<br/>Deployed?}
    
    InfraReady -->|No| ShowDeployError[❌ Show Deploy Error<br/>+ Logs]
    ShowDeployError --> Exit7([🚪 Exit 1])
    
    InfraReady -->|Yes| SetupAccess[🔌 Setup Access<br/>Port-forward + Browser]
    
    SetupAccess --> AccessReady{✅ Access<br/>Configured?}
    
    AccessReady -->|No| ShowAccessWarning[⚠️ Show Access Warning<br/>+ Manual instructions]
    ShowAccessWarning --> ShowSuccess
    
    AccessReady -->|Yes| ShowSuccess[🎉 Show Success Message<br/>Next steps + URLs]
    
    ShowSuccess --> Success([✅ Exit 0])
    
    style Start fill:#e1f5fe,color:#000000
    style Success fill:#e8f5e8,color:#000000
    style Exit1 fill:#ffebee,color:#000000
    style Exit2 fill:#ffebee,color:#000000
    style Exit3 fill:#ffebee,color:#000000
    style Exit4 fill:#ffebee,color:#000000
    style Exit5 fill:#fff3e0,color:#000000
    style Exit6 fill:#ffebee,color:#000000
    style Exit7 fill:#ffebee,color:#000000

```

## 🔄 Fluxo Secundário: Install Lab

### 8.3 Diagrama de Sequência - Install Lab

```mermaid

sequenceDiagram
    participant User as 👤 User
    participant CLI as 📱 CLI Command
    participant Handler as 🎯 Install Handler
    participant Orch as 🎼 Lab Orchestrator
    participant LabSvc as 🧪 Lab Service
    participant RepoSvc as 📚 Repository Service
    participant Validator as ✅ Lab Validator
    participant K8s as ☸️ K8s Service
    participant UI as 🎨 UI Service

    User->>CLI: girus create lab kubernetes-basics
    
    CLI->>Handler: Handle(ctx, InstallLabRequest)
    
    Handler->>Orch: InstallLab(ctx, request)
    
    Orch->>UI: ShowHeader("INSTALLING LAB")
    UI-->>User: 🎨 Display header
    
    Orch->>RepoSvc: FindLab(labID, repoURL)
    RepoSvc->>RepoSvc: Download index.yaml
    RepoSvc->>RepoSvc: Search for lab
    
    alt Lab Not Found
        RepoSvc-->>Orch: LabNotFoundError
        Orch->>UI: ShowError + Available labs
        UI-->>User: ❌ Error with suggestions
        Orch-->>Handler: Error
    else Lab Found
        RepoSvc-->>Orch: LabMetadata ✅
        
        Orch->>RepoSvc: DownloadLab(labURL)
        RepoSvc->>UI: ShowProgress("Downloading lab...")
        UI-->>User: 📊 Progress bar
        RepoSvc-->>Orch: LabContent
        
        Orch->>Validator: ValidateLab(content)
        Validator->>Validator: Check required fields
        Validator->>Validator: Validate resources
        Validator->>Validator: Check security rules
        
        alt Validation Failed
            Validator-->>Orch: ValidationErrors
            Orch->>UI: ShowValidationErrors
            UI-->>User: ❌ Validation errors
            Orch-->>Handler: Error
        else Validation Passed
            Validator-->>Orch: ValidationSuccess ✅
            
            Orch->>LabSvc: InstallLab(ctx, labSpec)
            LabSvc->>K8s: ApplyConfigMap(labTemplate)
            K8s->>UI: ShowProgress("Installing lab...")
            UI-->>User: 📊 Progress bar
            K8s-->>LabSvc: ConfigMap applied ✅
            
            LabSvc->>K8s: RestartBackend()
            K8s->>UI: ShowProgress("Restarting backend...")
            UI-->>User: 📊 Progress bar
            K8s->>K8s: Rollout restart deployment
            K8s->>K8s: Wait for ready
            K8s-->>LabSvc: Backend restarted ✅
            
            LabSvc-->>Orch: InstallResult
            Orch->>UI: ShowSuccess("Lab installed!")
            UI-->>User: ✅ Success + Lab info
            
            Orch-->>Handler: Success
            Handler-->>CLI: Success
            CLI-->>User: Exit code 0
        end
    end

```

### 8.4 Diagrama de Componentes - Arquitetura Geral

```mermaid

graph TB
    subgraph "🖥️ CLI Layer"
        CMD[Commands<br/>create, delete, list]
    end
    
    subgraph "🎯 Application Layer"
        HANDLER[Handlers<br/>Request processing]
        ORCH[Orchestrators<br/>Business workflows]
        DTO[DTOs<br/>Data transfer]
    end
    
    subgraph "🏗️ Domain Layer"
        CLUSTER[Cluster Domain<br/>Business logic]
        LAB[Lab Domain<br/>Validation rules]
        REPO[Repository Domain<br/>Index management]
        INFRA[Infrastructure Domain<br/>Prerequisites]
    end
    
    subgraph "🔧 Infrastructure Layer"
        K8S[Kubernetes<br/>Client & Operations]
        DOCKER[Container Engine<br/>Docker/Podman]
        HTTP[HTTP Client<br/>Repository access]
        FS[Filesystem<br/>File operations]
        EXEC[Command Executor<br/>Shell commands]
    end
    
    subgraph "🔄 Adapters Layer"
        CONFIG[Config Loader<br/>Multi-source]
        LOG[Logging<br/>Structured logs]
        I18N[Internationalization<br/>Messages]
        ERROR[Error Handler<br/>Recovery strategies]
        UI[UI Service<br/>Progress & display]
    end
    
    subgraph "📦 External Systems"
        KIND[Kind CLI]
        KUBECTL[Kubectl]
        REGISTRY[Container Registry]
        GITHUB[GitHub Repos]
    end
    
    CMD --> HANDLER
    HANDLER --> ORCH
    ORCH --> CLUSTER
    ORCH --> LAB
    ORCH --> REPO
    ORCH --> INFRA
    
    CLUSTER --> K8S
    LAB --> K8S
    REPO --> HTTP
    INFRA --> DOCKER
    INFRA --> EXEC
    
    K8S --> KIND
    K8S --> KUBECTL
    DOCKER --> REGISTRY
    HTTP --> GITHUB
    EXEC --> KIND
    EXEC --> KUBECTL
    
    HANDLER --> CONFIG
    HANDLER --> ERROR
    ORCH --> UI
    ORCH --> LOG
    ERROR --> I18N
    UI --> I18N
    
    style CMD fill:#e3f2fd,color:#000000
    style HANDLER fill:#f3e5f5,color:#000000
    style ORCH fill:#f3e5f5,color:#000000
    style CLUSTER fill:#e8f5e8,color:#000000
    style LAB fill:#e8f5e8,color:#000000
    style K8S fill:#fff3e0,color:#000000
    style DOCKER fill:#fff3e0,color:#000000
    style CONFIG fill:#f1f8e9,color:#000000
    style ERROR fill:#f1f8e9,color:#000000

```

### 8.5 Fluxo de Dependências

```mermaid

graph LR
    subgraph "📱 Presentation"
        CLI[CLI Commands]
    end
    
    subgraph "🎯 Application"
        HANDLERS[Handlers]
        ORCH[Orchestrators]
    end
    
    subgraph "🏗️ Domain"
        SERVICES[Domain Services]
        ENTITIES[Entities]
        REPOS[Repository Interfaces]
    end
    
    subgraph "🔧 Infrastructure"
        REPO_IMPL[Repository Implementations]
        EXT_CLIENTS[External Clients]
    end
    
    subgraph "🔄 Cross-Cutting"
        CONFIG[Configuration]
        LOGGING[Logging]
        ERRORS[Error Handling]
    end
    
    CLI --> HANDLERS
    HANDLERS --> ORCH
    ORCH --> SERVICES
    SERVICES --> ENTITIES
    SERVICES --> REPOS
    REPOS --> REPO_IMPL
    REPO_IMPL --> EXT_CLIENTS
    
    HANDLERS -.-> CONFIG
    ORCH -.-> LOGGING
    SERVICES -.-> ERRORS
    REPO_IMPL -.-> CONFIG
    
    style CLI fill:#e3f2fd,color:#000000
    style HANDLERS fill:#f3e5f5,color:#000000
    style SERVICES fill:#e8f5e8,color:#000000
    style REPO_IMPL fill:#fff3e0,color:#000000
    style CONFIG fill:#f1f8e9,color:#000000

```

## 🎯 Vantagens do Novo Fluxo

### 8.6 Comparação de Complexidade

| Aspecto | Arquitetura Atual | Nova Arquitetura |
|---------|------------------|------------------|
| **Pontos de decisão** | 20+ em uma função | 3-5 por componente |
| **Dependências** | Hardcoded | Injetadas via interfaces |
| **Testabilidade** | Impossível | Cada componente testável |
| **Reutilização** | Baixa | Alta (serviços modulares) |
| **Debugging** | Difícil (500 linhas) | Fácil (responsabilidade única) |
| **Manutenção** | Alto risco | Baixo risco (mudanças isoladas) |

### 8.7 Fluxo de Erros e Recovery

```mermaid

flowchart TD
    Error[❌ Error Occurs] --> Classify{🔍 Classify Error<br/>Code & Severity}
    
    Classify --> Log[📝 Log Error<br/>Structured logging]
    
    Log --> Recoverable{🔧 Is Recoverable?}
    
    Recoverable -->|Yes| FindStrategy[🎯 Find Recovery Strategy]
    FindStrategy --> TryRecover[🔄 Attempt Recovery]
    
    TryRecover --> RecoverSuccess{✅ Recovery<br/>Successful?}
    
    RecoverSuccess -->|Yes| ShowRecoverySuccess[🎉 Show Recovery Success]
    RecoverSuccess -->|No| ShowUserError
    
    Recoverable -->|No| ShowUserError[🎨 Show User-Friendly Error]
    
    ShowUserError --> Suggestions[💡 Show Suggestions<br/>Based on error type]
    
    Suggestions --> Localize[🌐 Localize Messages<br/>User's language]
    
    Localize --> Report[📊 Report Telemetry<br/>If enabled]
    
    Report --> End([🚪 Exit])
    
    ShowRecoverySuccess --> Continue([▶️ Continue Execution])
    
    style Error fill:#ffebee,color:#000000
    style ShowRecoverySuccess fill:#e8f5e8,color:#000000
    style Continue fill:#e8f5e8,color:#000000
    style End fill:#fff3e0,color:#000000

```

## 📊 Métricas do Novo Fluxo

### 8.8 Indicadores de Qualidade

| Métrica | Valor Atual | Meta Nova Arquitetura |
|---------|-------------|----------------------|
| **Complexidade Ciclomática** | 25+ | < 5 por função |
| **Linhas por função** | 500+ | < 50 |
| **Cobertura de testes** | 0% | > 80% |
| **Dependências por módulo** | 15+ | < 5 |
| **Tempo para adicionar feature** | 2-3 dias | 2-4 horas |
| **Tempo para debug** | 1-2 horas | 10-20 minutos |

### 8.9 Fluxo de Desenvolvimento

```mermaid

flowchart LR
    Req[📋 Requirement] --> Design[🎨 Domain Design]
    Design --> Interface[📝 Define Interfaces]
    Interface --> Test[🧪 Write Tests]
    Test --> Impl[⚡ Implement]
    Impl --> Integrate[🔗 Integration]
    Integrate --> Deploy[🚀 Deploy]
    
    Test -.-> Mock[🎭 Create Mocks]
    Mock -.-> Unit[🔬 Unit Tests]
    
    style Req fill:#e3f2fd,color:#000000
    style Test fill:#e8f5e8,color:#000000
    style Deploy fill:#f1f8e9,color:#000000

```

## 🔄 Próxima Etapa

[Plano de Implementação](./09-plano-implementacao.md) - Estratégia para implementar a nova arquitetura.
