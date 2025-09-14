# 6. Tratamento de Erros

## 🎯 Objetivo

Implementar um sistema de tratamento de erros consistente, estruturado e localizado que facilite debugging, recovery automático e experiência do usuário.

## 🏗️ Arquitetura de Erros

### 6.1 Hierarquia de Tipos de Erro

```txt
internal/errors/types.go
```

```go
type GirusError struct {
    Code        ErrorCode              `json:"code"`
    Message     string                 `json:"message"`
    Details     string                 `json:"details,omitempty"`
    Cause       error                  `json:"-"`
    Context     map[string]interface{} `json:"context,omitempty"`
    Timestamp   time.Time              `json:"timestamp"`
    Severity    ErrorSeverity          `json:"severity"`
    Recoverable bool                   `json:"recoverable"`
    
    // Localização
    MessageKey  string        `json:"-"`
    MessageArgs []interface{} `json:"-"`
    
    // Stack trace (apenas em modo debug)
    StackTrace []string `json:"stack_trace,omitempty"`
}

type ErrorCode string

const (
    // Infrastructure errors (1xxx)
    ErrContainerEngineNotFound    ErrorCode = "ERR_1001"
    ErrContainerEngineNotRunning  ErrorCode = "ERR_1002"
    ErrKindNotFound               ErrorCode = "ERR_1003"
    ErrKubectlNotFound            ErrorCode = "ERR_1004"
    ErrDockerDaemonNotRunning     ErrorCode = "ERR_1005"
    ErrPodmanServiceNotRunning    ErrorCode = "ERR_1006"
    ErrInsufficientResources      ErrorCode = "ERR_1007"
    ErrPortAlreadyInUse           ErrorCode = "ERR_1008"
    
    // Cluster errors (2xxx)
    ErrClusterAlreadyExists       ErrorCode = "ERR_2001"
    ErrClusterNotFound            ErrorCode = "ERR_2002"
    ErrClusterCreationFailed      ErrorCode = "ERR_2003"
    ErrClusterDeploymentFailed    ErrorCode = "ERR_2004"
    ErrClusterDeleteFailed        ErrorCode = "ERR_2005"
    ErrClusterNotReady            ErrorCode = "ERR_2006"
    ErrKubeConfigInvalid          ErrorCode = "ERR_2007"
    
    // Lab errors (3xxx)
    ErrLabFileNotFound            ErrorCode = "ERR_3001"
    ErrLabValidationFailed        ErrorCode = "ERR_3002"
    ErrLabInstallationFailed      ErrorCode = "ERR_3003"
    ErrLabAlreadyExists           ErrorCode = "ERR_3004"
    ErrLabUninstallFailed         ErrorCode = "ERR_3005"
    ErrLabTemplateInvalid         ErrorCode = "ERR_3006"
    ErrLabResourcesExceeded       ErrorCode = "ERR_3007"
    
    // Repository errors (4xxx)
    ErrRepositoryNotReachable     ErrorCode = "ERR_4001"
    ErrLabNotFoundInRepo          ErrorCode = "ERR_4002"
    ErrInvalidRepositoryFormat    ErrorCode = "ERR_4003"
    ErrRepositoryAuthFailed       ErrorCode = "ERR_4004"
    ErrRepositoryTimeout          ErrorCode = "ERR_4005"
    ErrIndexFileCorrupted         ErrorCode = "ERR_4006"
    
    // Configuration errors (5xxx)
    ErrConfigValidationFailed     ErrorCode = "ERR_5001"
    ErrConfigFileNotFound         ErrorCode = "ERR_5002"
    ErrConfigParsingFailed        ErrorCode = "ERR_5003"
    ErrInvalidConfigValue         ErrorCode = "ERR_5004"
    
    // Network errors (6xxx)
    ErrNetworkTimeout             ErrorCode = "ERR_6001"
    ErrDNSResolutionFailed        ErrorCode = "ERR_6002"
    ErrConnectionRefused          ErrorCode = "ERR_6003"
    ErrTLSCertificateInvalid      ErrorCode = "ERR_6004"
    
    // Permission errors (7xxx)
    ErrPermissionDenied           ErrorCode = "ERR_7001"
    ErrFilePermissionDenied       ErrorCode = "ERR_7002"
    ErrDockerPermissionDenied     ErrorCode = "ERR_7003"
    ErrKubernetesPermissionDenied ErrorCode = "ERR_7004"
    
    // General errors (9xxx)
    ErrInvalidInput               ErrorCode = "ERR_9001"
    ErrTimeout                    ErrorCode = "ERR_9002"
    ErrCancelled                  ErrorCode = "ERR_9003"
    ErrInternalError              ErrorCode = "ERR_9999"
)

type ErrorSeverity string

const (
    SeverityInfo     ErrorSeverity = "info"
    SeverityWarning  ErrorSeverity = "warning"
    SeverityError    ErrorSeverity = "error"
    SeverityCritical ErrorSeverity = "critical"
)

func (e *GirusError) Error() string {
    if e.Cause != nil {
        return fmt.Sprintf("[%s] %s (caused by: %s)", e.Code, e.Message, e.Cause.Error())
    }
    return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func (e *GirusError) Unwrap() error {
    return e.Cause
}

func (e *GirusError) WithContext(key string, value interface{}) *GirusError {
    if e.Context == nil {
        e.Context = make(map[string]interface{})
    }
    e.Context[key] = value
    return e
}

func (e *GirusError) WithDetails(details string) *GirusError {
    e.Details = details
    return e
}

func (e *GirusError) IsRecoverable() bool {
    return e.Recoverable
}

func (e *GirusError) GetSeverity() ErrorSeverity {
    return e.Severity
}

```

### 6.2 Factory Functions

```txt
internal/errors/factories.go
```

```go
func NewContainerEngineError(engine string, cause error) *GirusError {
    return &GirusError{
        Code:        ErrContainerEngineNotFound,
        Severity:    SeverityError,
        Recoverable: true,
        MessageKey:  "error.container_engine_not_found",
        MessageArgs: []interface{}{engine},
        Message:     fmt.Sprintf("Container engine '%s' not found or not running", engine),
        Cause:       cause,
        Context: map[string]interface{}{
            "engine": engine,
            "os":     runtime.GOOS,
            "arch":   runtime.GOARCH,
        },
        Timestamp: time.Now(),
    }
}

func NewClusterExistsError(name string) *GirusError {
    return &GirusError{
        Code:        ErrClusterAlreadyExists,
        Severity:    SeverityWarning,
        Recoverable: true,
        MessageKey:  "error.cluster_already_exists",
        MessageArgs: []interface{}{name},
        Message:     fmt.Sprintf("Cluster '%s' already exists", name),
        Context: map[string]interface{}{
            "cluster_name": name,
            "action":       "create",
        },
        Timestamp: time.Now(),
    }
}

func NewLabValidationError(validationErrors []ValidationError, labFile string) *GirusError {
    details := make([]string, len(validationErrors))
    for i, ve := range validationErrors {
        details[i] = ve.Error()
    }
    
    return &GirusError{
        Code:        ErrLabValidationFailed,
        Severity:    SeverityError,
        Recoverable: false,
        MessageKey:  "error.lab_validation_failed",
        Message:     "Lab file validation failed",
        Details:     strings.Join(details, "; "),
        Context: map[string]interface{}{
            "lab_file":           labFile,
            "validation_errors":  validationErrors,
            "error_count":        len(validationErrors),
        },
        Timestamp: time.Now(),
    }
}

func NewNetworkError(operation string, url string, cause error) *GirusError {
    code := ErrNetworkTimeout
    if strings.Contains(cause.Error(), "connection refused") {
        code = ErrConnectionRefused
    } else if strings.Contains(cause.Error(), "no such host") {
        code = ErrDNSResolutionFailed
    }
    
    return &GirusError{
        Code:        code,
        Severity:    SeverityError,
        Recoverable: true,
        MessageKey:  "error.network_operation_failed",
        MessageArgs: []interface{}{operation, url},
        Message:     fmt.Sprintf("Network operation '%s' failed for %s", operation, url),
        Cause:       cause,
        Context: map[string]interface{}{
            "operation": operation,
            "url":       url,
            "timeout":   "30s",
        },
        Timestamp: time.Now(),
    }
}

func NewPermissionError(operation string, path string, cause error) *GirusError {
    return &GirusError{
        Code:        ErrPermissionDenied,
        Severity:    SeverityCritical,
        Recoverable: false,
        MessageKey:  "error.permission_denied",
        MessageArgs: []interface{}{operation, path},
        Message:     fmt.Sprintf("Permission denied for operation '%s' on '%s'", operation, path),
        Cause:       cause,
        Context: map[string]interface{}{
            "operation": operation,
            "path":      path,
            "user":      os.Getenv("USER"),
            "uid":       os.Getuid(),
            "gid":       os.Getgid(),
        },
        Timestamp: time.Now(),
    }
}
```

### 6.3 Error Handler

```txt
internal/errors/handler.go
```

```go
type Handler struct {
    ui        ui.Service
    logger    logger.Logger
    i18n      i18n.Service
    reporter  Reporter
    recovery  RecoveryManager
    config    *config.Config
}

func NewHandler(
    ui ui.Service,
    logger logger.Logger,
    i18n i18n.Service,
    reporter Reporter,
    recovery RecoveryManager,
    config *config.Config,
) *Handler {
    return &Handler{
        ui:       ui,
        logger:   logger,
        i18n:     i18n,
        reporter: reporter,
        recovery: recovery,
        config:   config,
    }
}

func (h *Handler) Handle(err error) error {
    if err == nil {
        return nil
    }
    
    var girusErr *GirusError
    if !errors.As(err, &girusErr) {
        // Converter erro genérico para GirusError
        girusErr = h.wrapGenericError(err)
    }
    
    // 1. Log do erro
    h.logError(girusErr)
    
    // 2. Reportar telemetria (se habilitado)
    if h.reporter != nil {
        h.reporter.ReportError(girusErr)
    }
    
    // 3. Tentar recovery automático
    if girusErr.IsRecoverable() && h.recovery != nil {
        if recovered := h.attemptRecovery(girusErr); recovered {
            h.ui.ShowMessage(ui.MessageLevelSuccess, 
                h.i18n.Localize("error.recovery_successful"),
            )
            return nil
        }
    }
    
    // 4. Exibir erro para o usuário
    h.displayError(girusErr)
    
    return girusErr
}

func (h *Handler) logError(err *GirusError) {
    fields := map[string]interface{}{
        "error_code":   err.Code,
        "severity":     err.Severity,
        "recoverable":  err.Recoverable,
        "context":      err.Context,
        "timestamp":    err.Timestamp,
    }
    
    if err.Cause != nil {
        fields["cause"] = err.Cause.Error()
    }
    
    switch err.Severity {
    case SeverityInfo:
        h.logger.Info(err.Message, fields)
    case SeverityWarning:
        h.logger.Warn(err.Message, fields)
    case SeverityError:
        h.logger.Error(err.Message, fields)
    case SeverityCritical:
        h.logger.Error(err.Message, fields)
        // Também log do stack trace em modo debug
        if h.config.LogLevel == "debug" && len(err.StackTrace) > 0 {
            h.logger.Debug("stack trace", "trace", err.StackTrace)
        }
    }
}

func (h *Handler) displayError(err *GirusError) {
    // Localizar mensagem
    localizedMessage := h.i18n.Localize(err.MessageKey, err.MessageArgs...)
    if localizedMessage == "" {
        localizedMessage = err.Message
    }
    
    // Construir UI error
    uiError := &ui.ErrorDisplay{
        Title:       h.getErrorTitle(err.Code),
        Message:     localizedMessage,
        Details:     err.Details,
        Code:        string(err.Code),
        Severity:    string(err.Severity),
        Suggestions: h.getSuggestions(err),
        Context:     h.formatContext(err.Context),
        Timestamp:   err.Timestamp,
    }
    
    h.ui.ShowError(uiError)
}

func (h *Handler) getSuggestions(err *GirusError) []string {
    switch err.Code {
    case ErrContainerEngineNotFound:
        engine := err.Context["engine"].(string)
        return h.getContainerEngineInstallSuggestions(engine)
        
    case ErrClusterAlreadyExists:
        clusterName := err.Context["cluster_name"].(string)
        return []string{
            fmt.Sprintf("Use 'girus delete cluster %s' to remove the existing cluster", clusterName),
            "Or use a different cluster name with --name flag",
            "Use 'girus status cluster' to check current clusters",
        }
        
    case ErrLabValidationFailed:
        return []string{
            "Check the lab file format against the documentation",
            "Use 'girus validate lab <file>' for detailed validation",
            "See example lab files at: https://github.com/badtuxx/girus-labs",
        }
        
    case ErrPermissionDenied:
        return []string{
            "Check file/directory permissions",
            "Ensure your user has appropriate access",
            "For Docker: add user to docker group with 'sudo usermod -aG docker $USER'",
        }
        
    case ErrNetworkTimeout:
        return []string{
            "Check your internet connection",
            "Verify proxy settings if behind corporate firewall",
            "Try increasing timeout with --timeout flag",
        }
        
    default:
        return []string{
            "Check the logs for more details with --verbose",
            "Visit https://github.com/badtuxx/girus-cli/issues for known issues",
            "Consider reporting this issue if it persists",
        }
    }
}

func (h *Handler) getContainerEngineInstallSuggestions(engine string) []string {
    switch {
    case runtime.GOOS == "darwin" && engine == "docker":
        return []string{
            "Install Docker Desktop: https://www.docker.com/products/docker-desktop",
            "Or use Colima: brew install colima docker && colima start",
            "Ensure Docker Desktop is running",
        }
    case runtime.GOOS == "linux" && engine == "docker":
        return []string{
            "Install Docker: curl -fsSL https://get.docker.com | bash",
            "Start Docker service: sudo systemctl start docker",
            "Add user to docker group: sudo usermod -aG docker $USER",
        }
    case engine == "podman":
        return []string{
            "Install Podman: https://podman.io/getting-started/installation",
            "Start Podman service (if needed): systemctl --user start podman",
            "Initialize Podman machine (macOS): podman machine init && podman machine start",
        }
    default:
        return []string{
            fmt.Sprintf("Install and configure %s", engine),
            "Ensure the service is running",
            "Check installation documentation",
        }
    }
}
```

### 6.4 Recovery Strategies

```txt
internal/errors/recovery.go
```

```go
type RecoveryManager interface {
    CanRecover(err *GirusError) bool
    Recover(ctx context.Context, err *GirusError) error
}

type DefaultRecoveryManager struct {
    strategies []RecoveryStrategy
    ui         ui.Service
    logger     logger.Logger
}

type RecoveryStrategy interface {
    CanHandle(err *GirusError) bool
    Recover(ctx context.Context, err *GirusError) error
    GetDescription() string
}

func NewRecoveryManager(ui ui.Service, logger logger.Logger) *DefaultRecoveryManager {
    return &DefaultRecoveryManager{
        strategies: []RecoveryStrategy{
            NewContainerEngineRecovery(ui, logger),
            NewNetworkRetryRecovery(ui, logger),
            NewPermissionRecovery(ui, logger),
            NewClusterRecovery(ui, logger),
        },
        ui:     ui,
        logger: logger,
    }
}

func (rm *DefaultRecoveryManager) CanRecover(err *GirusError) bool {
    for _, strategy := range rm.strategies {
        if strategy.CanHandle(err) {
            return true
        }
    }
    return false
}

func (rm *DefaultRecoveryManager) Recover(ctx context.Context, err *GirusError) error {
    for _, strategy := range rm.strategies {
        if strategy.CanHandle(err) {
            rm.logger.Info("attempting recovery", 
                "strategy", strategy.GetDescription(),
                "error_code", err.Code,
            )
            
            if recoveryErr := strategy.Recover(ctx, err); recoveryErr == nil {
                rm.logger.Info("recovery successful", 
                    "strategy", strategy.GetDescription(),
                )
                return nil
            } else {
                rm.logger.Warn("recovery failed", 
                    "strategy", strategy.GetDescription(),
                    "error", recoveryErr,
                )
            }
        }
    }
    
    return fmt.Errorf("no recovery strategy available for error: %s", err.Code)
}

// Container Engine Recovery Strategy
type ContainerEngineRecovery struct {
    ui     ui.Service
    logger logger.Logger
}

func NewContainerEngineRecovery(ui ui.Service, logger logger.Logger) *ContainerEngineRecovery {
    return &ContainerEngineRecovery{ui: ui, logger: logger}
}

func (r *ContainerEngineRecovery) CanHandle(err *GirusError) bool {
    return err.Code == ErrContainerEngineNotFound || 
           err.Code == ErrContainerEngineNotRunning ||
           err.Code == ErrDockerDaemonNotRunning
}

func (r *ContainerEngineRecovery) Recover(ctx context.Context, err *GirusError) error {
    engine := err.Context["engine"].(string)
    
    switch err.Code {
    case ErrContainerEngineNotRunning, ErrDockerDaemonNotRunning:
        // Tentar iniciar o container engine
        return r.startContainerEngine(ctx, engine)
        
    case ErrContainerEngineNotFound:
        // Oferecer para instalar (se suportado)
        if r.canAutoInstall(engine) {
            shouldInstall, err := r.ui.Confirm(
                fmt.Sprintf("Would you like to install %s automatically?", engine),
                false,
            )
            if err != nil || !shouldInstall {
                return err
            }
            
            return r.installContainerEngine(ctx, engine)
        }
    }
    
    return fmt.Errorf("cannot recover from %s", err.Code)
}

func (r *ContainerEngineRecovery) GetDescription() string {
    return "Container Engine Recovery"
}

func (r *ContainerEngineRecovery) startContainerEngine(ctx context.Context, engine string) error {
    r.ui.ShowMessage(ui.MessageLevelInfo, 
        fmt.Sprintf("Attempting to start %s...", engine),
    )
    
    switch engine {
    case "docker":
        if runtime.GOOS == "darwin" {
            // macOS: tentar iniciar Docker Desktop
            return r.startDockerDesktop(ctx)
        } else {
            // Linux: iniciar serviço systemd
            return r.startSystemdService(ctx, "docker")
        }
    case "podman":
        return r.startSystemdService(ctx, "podman")
    }
    
    return fmt.Errorf("unsupported container engine: %s", engine)
}

```

### 6.5 Localização de Mensagens

```txt
internal/adapters/i18n/messages/pt.yaml
```

```yml
errors:
  container_engine_not_found: "Motor de container '%s' não encontrado ou não está rodando"
  cluster_already_exists: "Cluster '%s' já existe"
  lab_validation_failed: "Validação do laboratório falhou"
  network_operation_failed: "Operação de rede '%s' falhou para %s"
  permission_denied: "Permissão negada para operação '%s' em '%s'"
  recovery_successful: "Recuperação automática bem-sucedida"

suggestions:
  check_logs: "Verifique os logs para mais detalhes com --verbose"
  report_issue: "Considere reportar este problema se persistir"
  check_docs: "Consulte a documentação em https://github.com/badtuxx/girus-cli"

recovery:
  attempting: "Tentando recuperação automática..."
  successful: "Recuperação bem-sucedida"
  failed: "Recuperação falhou"
```

## 🧪 Testes de Error Handling

### 6.6 Testes Unitários

```txt
internal/errors/handler_test.go
```

```go
func TestErrorHandler_Handle(t *testing.T) {
    tests := []struct {
        name              string
        error             error
        expectRecovery    bool
        expectedUIDisplay bool
        expectedLog       string
    }{
        {
            name: "recoverable container engine error",
            error: NewContainerEngineError("docker", errors.New("not running")),
            expectRecovery: true,
            expectedUIDisplay: false, // Recovery successful
            expectedLog: "attempting recovery",
        },
        {
            name: "non-recoverable validation error",
            error: NewLabValidationError([]ValidationError{
                {Field: "name", Message: "required"},
            }, "test.yaml"),
            expectRecovery: false,
            expectedUIDisplay: true,
            expectedLog: "lab file validation failed",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mockUI := &ui.MockService{}
            mockLogger := &logger.MockLogger{}
            mockRecovery := &MockRecoveryManager{}
            
            handler := NewHandler(mockUI, mockLogger, nil, nil, mockRecovery, &config.Config{})
            
            // Setup mocks
            if tt.expectRecovery {
                mockRecovery.On("CanRecover", mock.Anything).Return(true)
                mockRecovery.On("Recover", mock.Anything, mock.Anything).Return(nil)
                mockUI.On("ShowMessage", ui.MessageLevelSuccess, mock.Anything)
            } else {
                mockUI.On("ShowError", mock.Anything)
            }
            
            err := handler.Handle(tt.error)
            
            if tt.expectRecovery {
                assert.NoError(t, err)
            } else {
                assert.Error(t, err)
            }
            
            mockUI.AssertExpectations(t)
            mockLogger.AssertExpectations(t)
        })
    }
}
```

## 📊 Benefícios do Sistema de Erros

| Aspecto | Benefício | Exemplo |
|---------|-----------|---------|
| **Consistência** | Todos os erros seguem mesmo padrão | Código, mensagem, contexto |
| **Debugging** | Contexto rico e estruturado | Stack trace, metadata |
| **Recovery** | Recuperação automática | Iniciar Docker automaticamente |
| **Localização** | Mensagens traduzidas | Português, inglês, espanhol |
| **Experiência** | Sugestões práticas | Como resolver o problema |

## 🔄 Próxima Etapa

[Nova Estrutura de Pastas](./07-estrutura-pastas.md) - Organização por domínios e camadas.
