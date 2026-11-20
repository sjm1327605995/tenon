package yoga

type CloneNodeFunc func(
	oldNode *Node,
	owner *Node,
	childIndex uint32,
) *Node

type Logger func(
	config *Config,
	node *Node,
	level LogLevel,
	format string,
	args ...any,
) int

type Config struct {
	cloneNodeCallback_    CloneNodeFunc
	logger_               Logger
	useWebDefaults_       bool
	printTree_            bool
	experimentalFeatures_ EnumBitset
	errata_               Errata
	context_              any
	pointScaleFactor_     float32
	version_              uint32
}

func ConfigNew() *Config {
	return NewConfig(DefaultLogger)
}

func NewConfig(logger Logger) *Config {
	return &Config{
		cloneNodeCallback_:    nil,
		logger_:               logger,
		useWebDefaults_:       false,
		printTree_:            false,
		experimentalFeatures_: NewEnumBitset(),
		errata_:               0,
		pointScaleFactor_:     1.0,
	}
}

func (config *Config) SetUseWebDefaults(useWebDefaults bool) {
	if config.useWebDefaults_ != useWebDefaults {
		config.useWebDefaults_ = useWebDefaults
		config.version_++
	}
}

func (config *Config) UseWebDefaults() bool {
	return config.useWebDefaults_
}

func (config *Config) SetShouldPrintTree(printTree bool) {
	config.printTree_ = printTree
}

func (config *Config) ShouldPrintTree() bool {
	return config.printTree_
}

func (config *Config) SetExperimentalFeatureEnabled(feature ExperimentalFeature, enabled bool) {
	if config.IsExperimentalFeatureEnabled(feature) == enabled {
		return
	}
	if enabled {
		config.experimentalFeatures_.Set(uint(feature))
	} else {
		config.experimentalFeatures_.Reset(uint(feature))
	}
	config.version_++
}

func (config *Config) IsExperimentalFeatureEnabled(feature ExperimentalFeature) bool {
	return config.experimentalFeatures_.Test(uint(feature))
}

func (config *Config) GetEnabledExperiments() EnumBitset {
	return config.experimentalFeatures_
}

// SetErrata
func (config *Config) SetErrata(errata Errata) {
	if config.errata_ != errata {
		config.errata_ = errata
		config.version_++
	}
}

// AddErrata
func (config *Config) AddErrata(errata Errata) {
	if config.errata_&errata != errata {
		config.errata_ |= errata
		config.version_++
	}
}

// RemoveErrata
func (config *Config) RemoveErrata(errata Errata) {
	if config.errata_&errata != 0 {
		config.errata_ &= ^errata
		config.version_++
	}
}

// GetErrata
func (config *Config) GetErrata() Errata {
	return config.errata_
}

// HasErrata
func (config *Config) HasErrata(errata Errata) bool {
	return config.errata_&errata != ErrataNone
}

// SetPointScaleFactor
func (config *Config) SetPointScaleFactor(pointScaleFactor float32) {
	if config.pointScaleFactor_ != pointScaleFactor {
		config.pointScaleFactor_ = pointScaleFactor
		config.version_++
	}
}

// GetVersion changes whenever a layout-affecting config value changes.
func (config *Config) GetVersion() uint32 { return config.version_ }

// GetPointScaleFactor
func (config *Config) GetPointScaleFactor() float32 {
	return config.pointScaleFactor_
}

// SetContext
func (config *Config) SetContext(context any) {
	config.context_ = context
}

// GetContext
func (config *Config) GetContext() any {
	return config.context_
}

// SetLogger
func (config *Config) SetLogger(logger Logger) {
	config.logger_ = logger
}

// log
func (config *Config) log(node *Node, level LogLevel, format string, args ...any) {
	if config.logger_ != nil {
		config.logger_(config, node, level, format, args...)
	}
}

// SetCloneNodeCallback
func (config *Config) SetCloneNodeCallback(callback CloneNodeFunc) {
	config.cloneNodeCallback_ = callback
}

// cloneNode
func (config *Config) cloneNode(oldNode *Node, owner *Node, childIndex uint32) *Node {
	var clone *Node
	if config.cloneNodeCallback_ != nil {
		clone = config.cloneNodeCallback_(oldNode, owner, childIndex)
	}
	if clone == nil {
		clone = oldNode.Clone()
	}
	return clone
}

var defaultConfig Config

func init() {
	defaultConfig = Config{
		logger_:           DefaultLogger,
		pointScaleFactor_: 1,
	}
}

func configUpdateInvalidatesLayout(oldConfig, newConfig *Config) bool {
	return oldConfig.GetErrata() != newConfig.GetErrata() ||
		oldConfig.GetEnabledExperiments() != newConfig.GetEnabledExperiments() ||
		oldConfig.GetPointScaleFactor() != newConfig.GetPointScaleFactor() ||
		oldConfig.UseWebDefaults() != newConfig.UseWebDefaults()
}
