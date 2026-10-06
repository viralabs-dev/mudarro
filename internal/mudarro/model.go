package mudarro

import "github.com/viralabs-dev/mudarro/internal/mudarro/model"

// Aliases preserve the orchestration API while adapters depend only on model.
type Command = model.Command
type Infrastructure = model.Infrastructure
type Database = model.Database
type Service = model.Service
type Suggestion = model.Suggestion
type Action = model.Action
