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
	config.useWebDefaults_ = useWebDefaults
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
	if enabled {
		config.experimentalFeatures_.Set(uint(feature))
	} else {
		config.experimentalFeatures_.Reset(uint(feature))
	}
}

func (config *Config) IsExperimentalFeatureEnabled(feature ExperimentalFeature) bool {
	return config.experimentalFeatures_.Test(uint(feature))
}

func (config *Config) GetEnabledExperiments() EnumBitset {
	return config.experimentalFeatures_
}

// SetErrata
func (config *Config) SetErrata(errata Errata) {
	config.errata_ = errata
}

// AddErrata
func (config *Config) AddErrata(errata Errata) {
	config.errata_ |= errata
}

// RemoveErrata
func (config *Config) RemoveErrata(errata Errata) {
	config.errata_ &= ^errata
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
	config.pointScaleFactor_ = pointScaleFactor
}

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
		clone = &Node{}
		*clone = *oldNode
		clone.setOwner(nil)
	}
	return clone
}

var defaultConfig Config

func init() {
	defaultConfig = Config{
		logger_: DefaultLogger,
	}
}

func configUpdateInvalidatesLayout(oldConfig, newConfig *Config) bool {
	return oldConfig.GetErrata() != newConfig.GetErrata() ||
		oldConfig.GetEnabledExperiments() != newConfig.GetEnabledExperiments() ||
		oldConfig.GetPointScaleFactor() != newConfig.GetPointScaleFactor() ||
		oldConfig.UseWebDefaults() != newConfig.UseWebDefaults()
}
