package step

import (
	"fmt"

	"github.com/bitrise-io/go-steputils/v2/cache"
	"github.com/bitrise-io/go-steputils/v2/stepconf"
	"github.com/bitrise-io/go-utils/v2/env"
	"github.com/bitrise-io/go-utils/v2/log"
	"github.com/bitrise-io/go-utils/v2/pathutil"
)

const (
	stepId = "save-npm-cache"

	// Cache key template
	// OS + Arch: guarantees unique cache per stack (node_modules can contain compiled binaries)
	// checksum: NPM or Yarn lockfiles
	// Note: it only considers the top-level lockfiles because there are other lockfiles inside `node_modules`,
	// otherwise the checksum would be different with and without a restored `node_modules` folder.
	key = `{{ .OS }}-{{ .Arch }}-npm-cache-{{ checksum "package-lock.json" "yarn.lock" }}`

	// Cached path
	path = "node_modules"
)

type Input struct {
	Verbose bool `env:"verbose,required"`
}

type SaveCacheStep struct {
	logger       log.Logger
	inputParser  stepconf.InputParser
	pathChecker  pathutil.PathChecker
	pathProvider pathutil.PathProvider
	pathModifier pathutil.PathModifier
	envRepo      env.Repository
}

func New(
	logger log.Logger,
	inputParser stepconf.InputParser,
	pathChecker pathutil.PathChecker,
	pathProvider pathutil.PathProvider,
	pathModifier pathutil.PathModifier,
	envRepo env.Repository,
) SaveCacheStep {
	return SaveCacheStep{
		logger:       logger,
		inputParser:  inputParser,
		pathChecker:  pathChecker,
		pathProvider: pathProvider,
		pathModifier: pathModifier,
		envRepo:      envRepo,
	}
}

func (step SaveCacheStep) Run() error {
	var input Input
	if err := step.inputParser.Parse(&input); err != nil {
		return fmt.Errorf("failed to parse inputs: %w", err)
	}
	stepconf.Print(input)
	step.logger.Println()
	step.logger.Printf("Cache key: %s", key)
	step.logger.Printf("Cache paths:")
	step.logger.Printf(path)
	step.logger.Println()

	step.logger.EnableDebugLog(input.Verbose)

	saver := cache.NewSaver(step.envRepo, step.logger, step.pathProvider, step.pathModifier, step.pathChecker)
	return saver.Save(cache.SaveCacheInput{
		StepId:      stepId,
		Verbose:     input.Verbose,
		Key:         key,
		Paths:       []string{path},
		IsKeyUnique: true,
	})
}
