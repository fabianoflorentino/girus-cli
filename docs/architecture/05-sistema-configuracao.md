# 5. Sistema de Configuração

## 🎯 Objetivo

Implementar um sistema de configuração robusto, hierárquico e type-safe que suporte múltiplas fontes de configuração com validação automática.

## 🏗️ Arquitetura de Configuração

### 5.1 Hierarquia de Configuração

```txt

Prioridade (maior para menor):
1. Flags da CLI (--cluster-name)
2. Variáveis de ambiente (GIRUS_CLUSTER_NAME)
3. Arquivo de configuração (~/.girus/config.yaml)
4. Valores padrão (defaults)

```

### 5.2 Estrutura Principal

```txt
internal/config/config.go
```

```go
type Config struct {
    // Core settings
    LogLevel  string `yaml:"log_level" env:"GIRUS_LOG_LEVEL" default:"info" validate:"oneof=debug info warn error"`
    LogFormat string `yaml:"log_format" env:"GIRUS_LOG_FORMAT" default:"text" validate:"oneof=text json"`
    Language  string `yaml:"language" env:"GIRUS_LANGUAGE" default:"pt" validate:"oneof=pt en es"`
    
    // Cluster configuration
    Cluster ClusterConfig `yaml:"cluster"`
    
    // Lab configuration
    Lab LabConfig `yaml:"lab"`
    
    // Repository configuration
    Repository RepositoryConfig `yaml:"repository"`
    
    // UI configuration
    UI UIConfig `yaml:"ui"`
    
    // Advanced settings
    Advanced AdvancedConfig `yaml:"advanced"`
}

type ClusterConfig struct {
    DefaultName         string        `yaml:"default_name" env:"GIRUS_CLUSTER_NAME" default:"girus" validate:"required,min=3,max=63"`
    ContainerEngine     string        `yaml:"container_engine" env:"GIRUS_CONTAINER_ENGINE" default:"docker" validate:"oneof=docker podman"`
    KubeConfig          string        `yaml:"kube_config" env:"KUBECONFIG" validate:"omitempty,file"`
    AutoPortForward     bool          `yaml:"auto_port_forward" env:"GIRUS_AUTO_PORT_FORWARD" default:"true"`
    AutoOpenBrowser     bool          `yaml:"auto_open_browser" env:"GIRUS_AUTO_OPEN_BROWSER" default:"true"`
    
    // Resource limits
    Resources ResourceConfig `yaml:"resources"`
    
    // Timeout configurations
    Timeouts TimeoutConfig `yaml:"timeouts"`
    
    // Default deployment files
    DefaultDeployments []string `yaml:"default_deployments"`
}

type ResourceConfig struct {
    DefaultCPU      string `yaml:"default_cpu" default:"2" validate:"required"`
    DefaultMemory   string `yaml:"default_memory" default:"4Gi" validate:"required"`
    DefaultStorage  string `yaml:"default_storage" default:"20Gi" validate:"required"`
    
    // Limits
    MaxCPU          string `yaml:"max_cpu" default:"8" validate:"required"`
    MaxMemory       string `yaml:"max_memory" default:"16Gi" validate:"required"`
    MaxStorage      string `yaml:"max_storage" default:"100Gi" validate:"required"`
}

type TimeoutConfig struct {
    ClusterCreate   time.Duration `yaml:"cluster_create" default:"10m" validate:"required"`
    PodReady        time.Duration `yaml:"pod_ready" default:"5m" validate:"required"`
    PortForward     time.Duration `yaml:"port_forward" default:"30s" validate:"required"`
    LabInstall      time.Duration `yaml:"lab_install" default:"3m" validate:"required"`
    HealthCheck     time.Duration `yaml:"health_check" default:"60s" validate:"required"`
}

type LabConfig struct {
    DefaultRepository string   `yaml:"default_repository" env:"GIRUS_DEFAULT_REPO" default:"https://raw.githubusercontent.com/badtuxx/girus-labs/main/index.yaml"`
    CacheDir         string   `yaml:"cache_dir" env:"GIRUS_CACHE_DIR" validate:"omitempty,dir"`
    CacheTTL         time.Duration `yaml:"cache_ttl" default:"1h" validate:"required"`
    AllowedSources   []string `yaml:"allowed_sources" default:"[\"file\", \"https\", \"http\"]"`
    
    // Validation settings
    Validation ValidationConfig `yaml:"validation"`
}

type ValidationConfig struct {
    StrictMode      bool     `yaml:"strict_mode" default:"false"`
    RequiredFields  []string `yaml:"required_fields"`
    MaxImageSize    string   `yaml:"max_image_size" default:"2Gi" validate:"required"`
    AllowedImages   []string `yaml:"allowed_images"`
    BlockedImages   []string `yaml:"blocked_images"`
}

type RepositoryConfig struct {
    IndexURLs           []string      `yaml:"index_urls"`
    CacheTimeout        time.Duration `yaml:"cache_timeout" default:"1h" validate:"required"`
    RetryAttempts       int          `yaml:"retry_attempts" default:"3" validate:"min=1,max=10"`
    ConnectionTimeout   time.Duration `yaml:"connection_timeout" default:"30s" validate:"required"`
    InsecureSkipVerify  bool         `yaml:"insecure_skip_verify" default:"false"`
    
    // Authentication (for private repositories)
    Auth AuthConfig `yaml:"auth"`
}

type AuthConfig struct {
    Username string `yaml:"username" env:"GIRUS_REPO_USERNAME"`
    Password string `yaml:"password" env:"GIRUS_REPO_PASSWORD"`
    Token    string `yaml:"token" env:"GIRUS_REPO_TOKEN"`
}

type UIConfig struct {
    ColorEnabled    bool   `yaml:"color_enabled" env:"GIRUS_COLOR" default:"true"`
    ProgressStyle   string `yaml:"progress_style" default:"bar" validate:"oneof=bar spinner dots"`
    VerboseDefault  bool   `yaml:"verbose_default" default:"false"`
    ConfirmActions  bool   `yaml:"confirm_actions" default:"true"`
    TableFormat     string `yaml:"table_format" default:"table" validate:"oneof=table json yaml"`
    
    // Internationalization
    DateFormat      string `yaml:"date_format" default:"2006-01-02 15:04:05"`
    Timezone        string `yaml:"timezone" default:"Local"`
}

type AdvancedConfig struct {
    // Performance tuning
    MaxConcurrentOps    int           `yaml:"max_concurrent_ops" default:"5" validate:"min=1,max=20"`
    RequestTimeout      time.Duration `yaml:"request_timeout" default:"2m" validate:"required"`
    
    // Debug settings
    EnableProfiling     bool   `yaml:"enable_profiling" default:"false"`
    EnableMetrics       bool   `yaml:"enable_metrics" default:"false"`
    MetricsPort         int    `yaml:"metrics_port" default:"9090" validate:"min=1024,max=65535"`
    
    // Feature flags
    FeatureFlags        map[string]bool `yaml:"feature_flags"`
}
```

### 5.3 Loader de Configuração

```txt
internal/config/loader.go
```

```go
type Loader struct {
    configPaths []string
    envPrefix   string
    validator   *validator.Validate
}

func NewLoader() *Loader {
    return &Loader{
        configPaths: []string{
            "./girus.yaml",
            "./girus.yml",
            "$HOME/.girus/config.yaml",
            "$HOME/.girus/config.yml",
            "$HOME/.config/girus/config.yaml",
            "/etc/girus/config.yaml",
        },
        envPrefix: "GIRUS_",
        validator: validator.New(),
    }
}

func (l *Loader) Load() (*Config, error) {
    config := &Config{}
    
    // 1. Aplicar valores padrão
    if err := l.applyDefaults(config); err != nil {
        return nil, fmt.Errorf("failed to apply defaults: %w", err)
    }
    
    // 2. Carregar arquivo de configuração
    if err := l.loadFromFile(config); err != nil {
        return nil, fmt.Errorf("failed to load config file: %w", err)
    }
    
    // 3. Aplicar variáveis de ambiente
    if err := l.loadFromEnv(config); err != nil {
        return nil, fmt.Errorf("failed to load env vars: %w", err)
    }
    
    // 4. Aplicar flags (será feito pelo comando CLI)
    
    // 5. Validar configuração final
    if err := l.validate(config); err != nil {
        return nil, fmt.Errorf("config validation failed: %w", err)
    }
    
    // 6. Processar configuração (expandir paths, etc.)
    if err := l.postProcess(config); err != nil {
        return nil, fmt.Errorf("config post-processing failed: %w", err)
    }
    
    return config, nil
}

func (l *Loader) applyDefaults(config *Config) error {
    // Usar reflection e tags para aplicar defaults
    return applyStructDefaults(config)
}

func (l *Loader) loadFromFile(config *Config) error {
    for _, path := range l.configPaths {
        expandedPath := expandPath(path)
        
        if _, err := os.Stat(expandedPath); os.IsNotExist(err) {
            continue
        }
        
        data, err := os.ReadFile(expandedPath)
        if err != nil {
            return fmt.Errorf("failed to read config file %s: %w", expandedPath, err)
        }
        
        if err := yaml.Unmarshal(data, config); err != nil {
            return fmt.Errorf("failed to parse config file %s: %w", expandedPath, err)
        }
        
        // Primeiro arquivo encontrado é usado
        break
    }
    
    return nil
}

func (l *Loader) loadFromEnv(config *Config) error {
    return applyEnvVars(config, l.envPrefix)
}

func (l *Loader) validate(config *Config) error {
    if err := l.validator.Struct(config); err != nil {
        return &ValidationError{
            Message: "configuration validation failed",
            Errors:  formatValidationErrors(err),
        }
    }
    
    // Validações customizadas
    if err := l.validateCustomRules(config); err != nil {
        return err
    }
    
    return nil
}

func (l *Loader) postProcess(config *Config) error {
    // Expandir paths relativos
    if config.Lab.CacheDir != "" {
        config.Lab.CacheDir = expandPath(config.Lab.CacheDir)
    }
    
    // Criar diretórios necessários
    if config.Lab.CacheDir != "" {
        if err := os.MkdirAll(config.Lab.CacheDir, 0755); err != nil {
            return fmt.Errorf("failed to create cache dir: %w", err)
        }
    }
    
    // Validar URLs de repositório
    for _, url := range config.Repository.IndexURLs {
        if err := validateURL(url); err != nil {
            return fmt.Errorf("invalid repository URL %s: %w", url, err)
        }
    }
    
    return nil
}
```

### 5.4 Validação Avançada

```txt
internal/config/validator.go
```

```go
type Validator struct {
    validate *validator.Validate
}

func NewValidator() *Validator {
    v := validator.New()
    
    // Registrar validações customizadas
    v.RegisterValidation("file", validateFileExists)
    v.RegisterValidation("dir", validateDirExists)
    v.RegisterValidation("semver", validateSemVer)
    v.RegisterValidation("k8sresource", validateK8sResource)
    
    return &Validator{validate: v}
}

func (v *Validator) ValidateConfig(config *Config) error {
    // Validação estrutural
    if err := v.validate.Struct(config); err != nil {
        return &ValidationError{
            Type:   "structural",
            Errors: v.formatValidationErrors(err),
        }
    }
    
    // Validações lógicas customizadas
    if err := v.validateBusinessRules(config); err != nil {
        return err
    }
    
    return nil
}

func (v *Validator) validateBusinessRules(config *Config) error {
    var errors []string
    
    // Validar recursos
    if err := v.validateResourceLimits(&config.Cluster.Resources); err != nil {
        errors = append(errors, err.Error())
    }
    
    // Validar timeouts
    if config.Cluster.Timeouts.ClusterCreate < time.Minute {
        errors = append(errors, "cluster create timeout must be at least 1 minute")
    }
    
    // Validar repositórios
    if len(config.Repository.IndexURLs) == 0 && config.Lab.DefaultRepository == "" {
        errors = append(errors, "at least one repository must be configured")
    }
    
    // Validar container engine
    if !isValidContainerEngine(config.Cluster.ContainerEngine) {
        errors = append(errors, fmt.Sprintf("unsupported container engine: %s", config.Cluster.ContainerEngine))
    }
    
    if len(errors) > 0 {
        return &ValidationError{
            Type:   "business_rules",
            Errors: errors,
        }
    }
    
    return nil
}

func validateK8sResource(fl validator.FieldLevel) bool {
    resource := fl.Field().String()
    
    // Validar formato de recursos Kubernetes (ex: "2", "500m", "2Gi")
    return regexp.MustCompile(`^(\d+(\.\d+)?[mMkKgGtTpPeE]?i?)$`).MatchString(resource)
}

type ValidationError struct {
    Type   string
    Errors []string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation failed (%s): %s", e.Type, strings.Join(e.Errors, ", "))
}
```

### 5.5 Configuração via Flags

```txt
cmd/root.go
```

```go
func initConfig() {
    // Carregar configuração base
    config, err := config.NewLoader().Load()
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }
    
    // Aplicar flags sobre a configuração
    applyFlagsToConfig(config)
    
    // Armazenar globalmente
    app.SetConfig(config)
}

func applyFlagsToConfig(config *Config) {
    // Cluster flags
    if clusterName != "" {
        config.Cluster.DefaultName = clusterName
    }
    if containerEngine != "" {
        config.Cluster.ContainerEngine = containerEngine
    }
    if verboseMode {
        config.LogLevel = "debug"
        config.UI.VerboseDefault = true
    }
    
    // Lab flags
    if labRepository != "" {
        config.Lab.DefaultRepository = labRepository
    }
    
    // UI flags
    if noColor {
        config.UI.ColorEnabled = false
    }
}

// Flags específicas por comando
func init() {
    createClusterCmd.Flags().StringVar(&clusterName, "name", "", "Nome do cluster (sobrescreve configuração)")
    createClusterCmd.Flags().StringVar(&containerEngine, "engine", "", "Container engine (docker|podman)")
    createClusterCmd.Flags().BoolVar(&verboseMode, "verbose", false, "Modo verbose")
    createClusterCmd.Flags().BoolVar(&skipPortForward, "skip-port-forward", false, "Pular configuração de port-forward")
}
```

## 📁 Exemplo de Arquivo de Configuração

### 5.6 ~/.girus/config.yaml

```yaml
# Configuração principal do GIRUS CLI
log_level: info
log_format: text
language: pt

cluster:
  default_name: girus
  container_engine: docker
  auto_port_forward: true
  auto_open_browser: true
  
  resources:
    default_cpu: "2"
    default_memory: "4Gi"
    default_storage: "20Gi"
    max_cpu: "8"
    max_memory: "16Gi"
    max_storage: "100Gi"
  
  timeouts:
    cluster_create: 10m
    pod_ready: 5m
    port_forward: 30s
    lab_install: 3m
  
  default_deployments:
    - "girus-infrastructure.yaml"
    - "girus-templates.yaml"

lab:
  default_repository: "https://raw.githubusercontent.com/badtuxx/girus-labs/main/index.yaml"
  cache_dir: "~/.girus/cache"
  cache_ttl: 1h
  allowed_sources:
    - "file"
    - "https"
  
  validation:
    strict_mode: false
    max_image_size: "2Gi"
    required_fields:
      - "metadata.name"
      - "spec.environment"

repository:
  index_urls:
    - "https://raw.githubusercontent.com/badtuxx/girus-labs/main/index.yaml"
  cache_timeout: 1h
  retry_attempts: 3
  connection_timeout: 30s
  insecure_skip_verify: false

ui:
  color_enabled: true
  progress_style: bar
  verbose_default: false
  confirm_actions: true
  table_format: table
  date_format: "2006-01-02 15:04:05"

advanced:
  max_concurrent_ops: 5
  request_timeout: 2m
  enable_profiling: false
  enable_metrics: false
  feature_flags:
    experimental_k3s: false
    helm_support: true
```

## 🧪 Testes de Configuração

### 5.7 Testes Unitários

```txt
internal/config/loader_test.go
```

```go
func TestConfigLoader_Load(t *testing.T) {
    tests := []struct {
        name           string
        configFile     string
        envVars        map[string]string
        expectedConfig *Config
        expectedError  string
    }{
        {
            name: "default configuration",
            expectedConfig: &Config{
                LogLevel:  "info",
                Language:  "pt",
                Cluster: ClusterConfig{
                    DefaultName:     "girus",
                    ContainerEngine: "docker",
                },
            },
        },
        {
            name: "override with env vars",
            envVars: map[string]string{
                "GIRUS_LOG_LEVEL":        "debug",
                "GIRUS_CLUSTER_NAME":     "custom-cluster",
                "GIRUS_CONTAINER_ENGINE": "podman",
            },
            expectedConfig: &Config{
                LogLevel: "debug",
                Cluster: ClusterConfig{
                    DefaultName:     "custom-cluster",
                    ContainerEngine: "podman",
                },
            },
        },
        {
            name: "invalid container engine",
            envVars: map[string]string{
                "GIRUS_CONTAINER_ENGINE": "invalid",
            },
            expectedError: "unsupported container engine",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Setup environment
            for key, value := range tt.envVars {
                t.Setenv(key, value)
            }
            
            loader := NewLoader()
            config, err := loader.Load()
            
            if tt.expectedError != "" {
                assert.Error(t, err)
                assert.Contains(t, err.Error(), tt.expectedError)
                return
            }
            
            assert.NoError(t, err)
            assert.Equal(t, tt.expectedConfig.LogLevel, config.LogLevel)
            assert.Equal(t, tt.expectedConfig.Cluster.DefaultName, config.Cluster.DefaultName)
        })
    }
}
```

## 📊 Benefícios do Sistema de Configuração

| Aspecto | Benefício | Exemplo |
|---------|-----------|---------|
| **Flexibilidade** | Múltiplas fontes | Arquivo, env vars, flags |
| **Validação** | Type-safe e regras de negócio | `validate:"oneof=docker podman"` |
| **Defaults** | Configuração zero | Funciona out-of-the-box |
| **Hierarquia** | Precedência clara | Flags > Env > File > Default |
| **Documentação** | Self-documenting | Tags e comentários |

## 🔄 Próxima Etapa

[Tratamento de Erros](./06-tratamento-erros.md) - Sistema robusto de error handling.
