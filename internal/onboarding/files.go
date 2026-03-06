package onboarding

import "github.com/cfpperche/vibeforge/internal/i18n"

type FillLevel int

const (
	FilledByScaffold FillLevel = iota // scaffold fills everything
	FilledPartial                     // scaffold fills structure, agent completes
	FilledByAgent                     // scaffold creates empty, agent fills
)

type MDFile struct {
	Path          string
	FillLevel     FillLevel
	Description   string
	ScaffoldFills []string
	AgentFills    []string
}

type Category struct {
	Name  string
	Files []string
}

// Categories returns the list with translated names.
// Rebuilt each call so language switches take effect.
var Categories []Category

// Rebuild must be called after i18n.Init() to populate translated strings.
func Rebuild() {
	RebuildCategories()
	RebuildFiles()
}

func RebuildCategories() {
	Categories = []Category{
		{i18n.T("onboarding.cat.product"), []string{"docs/PRODUCT_BRIEF.md", "docs/PERSONA.md", "docs/VIRAL_LOOP.md"}},
		{i18n.T("onboarding.cat.root"), []string{"README.md", "CHANGELOG.md", "LICENSE", "CODE_OF_CONDUCT.md", ".editorconfig"}},
		{i18n.T("onboarding.cat.claude"), []string{"CLAUDE.md", ".claude/settings.json", ".claude/hooks/", ".claude/commands/"}},
		{i18n.T("onboarding.cat.docs"), []string{"docs/CONTEXT.md", "docs/ROADMAP.md", "docs/ARCHITECTURE.md", "docs/GLOSSARY.md", "docs/TESTING.md", "docs/DEPLOYMENT.md"}},
		{i18n.T("onboarding.cat.requirements"), []string{"docs/requirements/SRS.md", "docs/requirements/RF.md", "docs/requirements/RNF.md", "docs/requirements/USER_STORIES.md", "docs/requirements/USE_CASES.md"}},
		{i18n.T("onboarding.cat.decisions"), []string{"docs/adr/0001-stack.md"}},
		{i18n.T("onboarding.cat.github"), []string{".github/workflows/ci.yml", ".github/workflows/release.yml", ".github/dependabot.yml", ".github/PULL_REQUEST_TEMPLATE.md", ".github/ISSUE_TEMPLATE/bug_report.md", ".github/ISSUE_TEMPLATE/feature_request.md"}},
		{i18n.T("onboarding.cat.quality"), []string{".pre-commit-config.yaml", "CONTRIBUTING.md", "SECURITY.md"}},
		{i18n.T("onboarding.cat.scripts"), []string{"scripts/setup.sh"}},
	}
}

// FileMap provides fast lookup by path.
var FileMap map[string]MDFile

// Files is the full list of file metadata.
var Files []MDFile

func RebuildFiles() {
	Files = []MDFile{
		// --- Product ---
		{
			Path: "docs/PRODUCT_BRIEF.md", FillLevel: FilledByScaffold,
			Description: i18n.T("onboarding.file.product_brief.desc"),
			ScaffoldFills: []string{
				i18n.T("onboarding.file.product_brief.scaffold.1"),
				i18n.T("onboarding.file.product_brief.scaffold.2"),
				i18n.T("onboarding.file.product_brief.scaffold.3"),
				i18n.T("onboarding.file.product_brief.scaffold.4"),
				i18n.T("onboarding.file.product_brief.scaffold.5"),
				i18n.T("onboarding.file.product_brief.scaffold.6"),
				i18n.T("onboarding.file.product_brief.scaffold.7"),
			},
			AgentFills: []string{
				i18n.T("onboarding.file.product_brief.agent.1"),
				i18n.T("onboarding.file.product_brief.agent.2"),
				i18n.T("onboarding.file.product_brief.agent.3"),
				i18n.T("onboarding.file.product_brief.agent.4"),
			},
		},
		{
			Path: "docs/PERSONA.md", FillLevel: FilledPartial,
			Description: i18n.T("onboarding.file.persona.desc"),
			ScaffoldFills: []string{
				i18n.T("onboarding.file.persona.scaffold.1"),
				i18n.T("onboarding.file.persona.scaffold.2"),
				i18n.T("onboarding.file.persona.scaffold.3"),
				i18n.T("onboarding.file.persona.scaffold.4"),
			},
			AgentFills: []string{
				i18n.T("onboarding.file.persona.agent.1"),
				i18n.T("onboarding.file.persona.agent.2"),
				i18n.T("onboarding.file.persona.agent.3"),
			},
		},
		{
			Path: "docs/VIRAL_LOOP.md", FillLevel: FilledPartial,
			Description: i18n.T("onboarding.file.viral_loop.desc"),
			ScaffoldFills: []string{
				i18n.T("onboarding.file.viral_loop.scaffold.1"),
				i18n.T("onboarding.file.viral_loop.scaffold.2"),
				i18n.T("onboarding.file.viral_loop.scaffold.3"),
				i18n.T("onboarding.file.viral_loop.scaffold.4"),
			},
			AgentFills: []string{
				i18n.T("onboarding.file.viral_loop.agent.1"),
				i18n.T("onboarding.file.viral_loop.agent.2"),
				i18n.T("onboarding.file.viral_loop.agent.3"),
			},
		},

		// --- Project root ---
		{
			Path: "README.md", FillLevel: FilledByScaffold,
			Description: i18n.T("onboarding.file.readme.desc"),
			ScaffoldFills: []string{
				i18n.T("onboarding.file.readme.scaffold.1"),
				i18n.T("onboarding.file.readme.scaffold.2"),
				i18n.T("onboarding.file.readme.scaffold.3"),
				i18n.T("onboarding.file.readme.scaffold.4"),
				i18n.T("onboarding.file.readme.scaffold.5"),
				i18n.T("onboarding.file.readme.scaffold.6"),
				i18n.T("onboarding.file.readme.scaffold.7"),
			},
			AgentFills: []string{
				i18n.T("onboarding.file.readme.agent.1"),
				i18n.T("onboarding.file.readme.agent.2"),
				i18n.T("onboarding.file.readme.agent.3"),
				i18n.T("onboarding.file.readme.agent.4"),
			},
		},
		{
			Path: "CHANGELOG.md", FillLevel: FilledPartial,
			Description: i18n.T("onboarding.file.changelog.desc"),
			ScaffoldFills: []string{
				i18n.T("onboarding.file.changelog.scaffold.1"),
				i18n.T("onboarding.file.changelog.scaffold.2"),
				i18n.T("onboarding.file.changelog.scaffold.3"),
			},
			AgentFills: []string{
				i18n.T("onboarding.file.changelog.agent.1"),
				i18n.T("onboarding.file.changelog.agent.2"),
				i18n.T("onboarding.file.changelog.agent.3"),
			},
		},
		{
			Path: "LICENSE", FillLevel: FilledByScaffold,
			Description: i18n.T("onboarding.file.license.desc"),
			ScaffoldFills: []string{
				i18n.T("onboarding.file.license.scaffold.1"),
				i18n.T("onboarding.file.license.scaffold.2"),
			},
			AgentFills: []string{},
		},
		{
			Path: "CODE_OF_CONDUCT.md", FillLevel: FilledByScaffold,
			Description: i18n.T("onboarding.file.coc.desc"),
			ScaffoldFills: []string{
				i18n.T("onboarding.file.coc.scaffold.1"),
				i18n.T("onboarding.file.coc.scaffold.2"),
			},
			AgentFills: []string{},
		},
		{
			Path: ".editorconfig", FillLevel: FilledByScaffold,
			Description: i18n.T("onboarding.file.editorconfig.desc"),
			ScaffoldFills: []string{
				i18n.T("onboarding.file.editorconfig.scaffold.1"),
				i18n.T("onboarding.file.editorconfig.scaffold.2"),
				i18n.T("onboarding.file.editorconfig.scaffold.3"),
			},
			AgentFills: []string{
				i18n.T("onboarding.file.editorconfig.agent.1"),
			},
		},

		// --- Claude Code ---
		{
			Path: "CLAUDE.md", FillLevel: FilledByScaffold,
			Description: i18n.T("onboarding.file.claude_md.desc"),
			ScaffoldFills: []string{
				i18n.T("onboarding.file.claude_md.scaffold.1"),
				i18n.T("onboarding.file.claude_md.scaffold.2"),
				i18n.T("onboarding.file.claude_md.scaffold.3"),
				i18n.T("onboarding.file.claude_md.scaffold.4"),
				i18n.T("onboarding.file.claude_md.scaffold.5"),
			},
			AgentFills: []string{
				i18n.T("onboarding.file.claude_md.agent.1"),
				i18n.T("onboarding.file.claude_md.agent.2"),
				i18n.T("onboarding.file.claude_md.agent.3"),
				i18n.T("onboarding.file.claude_md.agent.4"),
			},
		},
		{
			Path: ".claude/settings.json", FillLevel: FilledByScaffold,
			Description: i18n.T("onboarding.file.settings.desc"),
			ScaffoldFills: []string{
				i18n.T("onboarding.file.settings.scaffold.1"),
				i18n.T("onboarding.file.settings.scaffold.2"),
				i18n.T("onboarding.file.settings.scaffold.3"),
			},
			AgentFills: []string{
				i18n.T("onboarding.file.settings.agent.1"),
			},
		},
		{
			Path: ".claude/hooks/", FillLevel: FilledByScaffold,
			Description: i18n.T("onboarding.file.hooks.desc"),
			ScaffoldFills: []string{
				i18n.T("onboarding.file.hooks.scaffold.1"),
				i18n.T("onboarding.file.hooks.scaffold.2"),
			},
			AgentFills: []string{
				i18n.T("onboarding.file.hooks.agent.1"),
				i18n.T("onboarding.file.hooks.agent.2"),
			},
		},
		{
			Path: ".claude/commands/", FillLevel: FilledByScaffold,
			Description: i18n.T("onboarding.file.commands.desc"),
			ScaffoldFills: []string{
				i18n.T("onboarding.file.commands.scaffold.1"),
				i18n.T("onboarding.file.commands.scaffold.2"),
				i18n.T("onboarding.file.commands.scaffold.3"),
				i18n.T("onboarding.file.commands.scaffold.4"),
			},
			AgentFills: []string{
				i18n.T("onboarding.file.commands.agent.1"),
			},
		},

		// --- Documentation ---
		{
			Path: "docs/CONTEXT.md", FillLevel: FilledByScaffold,
			Description: i18n.T("onboarding.file.context.desc"),
			ScaffoldFills: []string{
				i18n.T("onboarding.file.context.scaffold.1"),
				i18n.T("onboarding.file.context.scaffold.2"),
				i18n.T("onboarding.file.context.scaffold.3"),
				i18n.T("onboarding.file.context.scaffold.4"),
				i18n.T("onboarding.file.context.scaffold.5"),
			},
			AgentFills: []string{
				i18n.T("onboarding.file.context.agent.1"),
				i18n.T("onboarding.file.context.agent.2"),
				i18n.T("onboarding.file.context.agent.3"),
				i18n.T("onboarding.file.context.agent.4"),
			},
		},
		{
			Path: "docs/ROADMAP.md", FillLevel: FilledPartial,
			Description: i18n.T("onboarding.file.roadmap.desc"),
			ScaffoldFills: []string{
				i18n.T("onboarding.file.roadmap.scaffold.1"),
				i18n.T("onboarding.file.roadmap.scaffold.2"),
				i18n.T("onboarding.file.roadmap.scaffold.3"),
			},
			AgentFills: []string{
				i18n.T("onboarding.file.roadmap.agent.1"),
				i18n.T("onboarding.file.roadmap.agent.2"),
				i18n.T("onboarding.file.roadmap.agent.3"),
			},
		},
		{
			Path: "docs/ARCHITECTURE.md", FillLevel: FilledByAgent,
			Description: i18n.T("onboarding.file.architecture.desc"),
			ScaffoldFills: []string{
				i18n.T("onboarding.file.architecture.scaffold.1"),
				i18n.T("onboarding.file.architecture.scaffold.2"),
			},
			AgentFills: []string{
				i18n.T("onboarding.file.architecture.agent.1"),
				i18n.T("onboarding.file.architecture.agent.2"),
				i18n.T("onboarding.file.architecture.agent.3"),
				i18n.T("onboarding.file.architecture.agent.4"),
			},
		},
		{
			Path: "docs/GLOSSARY.md", FillLevel: FilledPartial,
			Description: i18n.T("onboarding.file.glossary.desc"),
			ScaffoldFills: []string{
				i18n.T("onboarding.file.glossary.scaffold.1"),
				i18n.T("onboarding.file.glossary.scaffold.2"),
			},
			AgentFills: []string{
				i18n.T("onboarding.file.glossary.agent.1"),
				i18n.T("onboarding.file.glossary.agent.2"),
				i18n.T("onboarding.file.glossary.agent.3"),
				i18n.T("onboarding.file.glossary.agent.4"),
			},
		},
		{
			Path: "docs/TESTING.md", FillLevel: FilledPartial,
			Description: i18n.T("onboarding.file.testing.desc"),
			ScaffoldFills: []string{
				i18n.T("onboarding.file.testing.scaffold.1"),
				i18n.T("onboarding.file.testing.scaffold.2"),
				i18n.T("onboarding.file.testing.scaffold.3"),
				i18n.T("onboarding.file.testing.scaffold.4"),
			},
			AgentFills: []string{
				i18n.T("onboarding.file.testing.agent.1"),
				i18n.T("onboarding.file.testing.agent.2"),
				i18n.T("onboarding.file.testing.agent.3"),
				i18n.T("onboarding.file.testing.agent.4"),
			},
		},
		{
			Path: "docs/DEPLOYMENT.md", FillLevel: FilledPartial,
			Description: i18n.T("onboarding.file.deployment.desc"),
			ScaffoldFills: []string{
				i18n.T("onboarding.file.deployment.scaffold.1"),
				i18n.T("onboarding.file.deployment.scaffold.2"),
				i18n.T("onboarding.file.deployment.scaffold.3"),
				i18n.T("onboarding.file.deployment.scaffold.4"),
			},
			AgentFills: []string{
				i18n.T("onboarding.file.deployment.agent.1"),
				i18n.T("onboarding.file.deployment.agent.2"),
				i18n.T("onboarding.file.deployment.agent.3"),
				i18n.T("onboarding.file.deployment.agent.4"),
			},
		},

		// --- Requirements ---
		{
			Path: "docs/requirements/SRS.md", FillLevel: FilledByAgent,
			Description: i18n.T("onboarding.file.srs.desc"),
			ScaffoldFills: []string{
				i18n.T("onboarding.file.srs.scaffold.1"),
				i18n.T("onboarding.file.srs.scaffold.2"),
			},
			AgentFills: []string{
				i18n.T("onboarding.file.srs.agent.1"),
				i18n.T("onboarding.file.srs.agent.2"),
				i18n.T("onboarding.file.srs.agent.3"),
			},
		},
		{
			Path: "docs/requirements/RF.md", FillLevel: FilledPartial,
			Description: i18n.T("onboarding.file.rf.desc"),
			ScaffoldFills: []string{
				i18n.T("onboarding.file.rf.scaffold.1"),
				i18n.T("onboarding.file.rf.scaffold.2"),
			},
			AgentFills: []string{
				i18n.T("onboarding.file.rf.agent.1"),
				i18n.T("onboarding.file.rf.agent.2"),
				i18n.T("onboarding.file.rf.agent.3"),
				i18n.T("onboarding.file.rf.agent.4"),
			},
		},
		{
			Path: "docs/requirements/RNF.md", FillLevel: FilledPartial,
			Description: i18n.T("onboarding.file.rnf.desc"),
			ScaffoldFills: []string{
				i18n.T("onboarding.file.rnf.scaffold.1"),
				i18n.T("onboarding.file.rnf.scaffold.2"),
			},
			AgentFills: []string{
				i18n.T("onboarding.file.rnf.agent.1"),
				i18n.T("onboarding.file.rnf.agent.2"),
				i18n.T("onboarding.file.rnf.agent.3"),
			},
		},
		{
			Path: "docs/requirements/USER_STORIES.md", FillLevel: FilledByAgent,
			Description: i18n.T("onboarding.file.user_stories.desc"),
			ScaffoldFills: []string{
				i18n.T("onboarding.file.user_stories.scaffold.1"),
				i18n.T("onboarding.file.user_stories.scaffold.2"),
				i18n.T("onboarding.file.user_stories.scaffold.3"),
			},
			AgentFills: []string{
				i18n.T("onboarding.file.user_stories.agent.1"),
				i18n.T("onboarding.file.user_stories.agent.2"),
				i18n.T("onboarding.file.user_stories.agent.3"),
				i18n.T("onboarding.file.user_stories.agent.4"),
			},
		},
		{
			Path: "docs/requirements/USE_CASES.md", FillLevel: FilledByAgent,
			Description: i18n.T("onboarding.file.use_cases.desc"),
			ScaffoldFills: []string{
				i18n.T("onboarding.file.use_cases.scaffold.1"),
				i18n.T("onboarding.file.use_cases.scaffold.2"),
			},
			AgentFills: []string{
				i18n.T("onboarding.file.use_cases.agent.1"),
				i18n.T("onboarding.file.use_cases.agent.2"),
				i18n.T("onboarding.file.use_cases.agent.3"),
			},
		},

		// --- Decisions ---
		{
			Path: "docs/adr/0001-stack.md", FillLevel: FilledPartial,
			Description: i18n.T("onboarding.file.adr.desc"),
			ScaffoldFills: []string{
				i18n.T("onboarding.file.adr.scaffold.1"),
				i18n.T("onboarding.file.adr.scaffold.2"),
				i18n.T("onboarding.file.adr.scaffold.3"),
			},
			AgentFills: []string{
				i18n.T("onboarding.file.adr.agent.1"),
				i18n.T("onboarding.file.adr.agent.2"),
				i18n.T("onboarding.file.adr.agent.3"),
				i18n.T("onboarding.file.adr.agent.4"),
			},
		},

		// --- GitHub ---
		{
			Path: ".github/workflows/ci.yml", FillLevel: FilledByScaffold,
			Description: i18n.T("onboarding.file.ci.desc"),
			ScaffoldFills: []string{
				i18n.T("onboarding.file.ci.scaffold.1"),
				i18n.T("onboarding.file.ci.scaffold.2"),
				i18n.T("onboarding.file.ci.scaffold.3"),
			},
			AgentFills: []string{
				i18n.T("onboarding.file.ci.agent.1"),
				i18n.T("onboarding.file.ci.agent.2"),
				i18n.T("onboarding.file.ci.agent.3"),
			},
		},
		{
			Path: ".github/workflows/release.yml", FillLevel: FilledByScaffold,
			Description: i18n.T("onboarding.file.release.desc"),
			ScaffoldFills: []string{
				i18n.T("onboarding.file.release.scaffold.1"),
				i18n.T("onboarding.file.release.scaffold.2"),
				i18n.T("onboarding.file.release.scaffold.3"),
			},
			AgentFills: []string{
				i18n.T("onboarding.file.release.agent.1"),
				i18n.T("onboarding.file.release.agent.2"),
				i18n.T("onboarding.file.release.agent.3"),
			},
		},
		{
			Path: ".github/dependabot.yml", FillLevel: FilledByScaffold,
			Description: i18n.T("onboarding.file.dependabot.desc"),
			ScaffoldFills: []string{
				i18n.T("onboarding.file.dependabot.scaffold.1"),
				i18n.T("onboarding.file.dependabot.scaffold.2"),
			},
			AgentFills: []string{},
		},
		{
			Path: ".github/PULL_REQUEST_TEMPLATE.md", FillLevel: FilledByScaffold,
			Description: i18n.T("onboarding.file.pr_template.desc"),
			ScaffoldFills: []string{
				i18n.T("onboarding.file.pr_template.scaffold.1"),
				i18n.T("onboarding.file.pr_template.scaffold.2"),
				i18n.T("onboarding.file.pr_template.scaffold.3"),
				i18n.T("onboarding.file.pr_template.scaffold.4"),
			},
			AgentFills: []string{
				i18n.T("onboarding.file.pr_template.agent.1"),
			},
		},
		{
			Path: ".github/ISSUE_TEMPLATE/bug_report.md", FillLevel: FilledByScaffold,
			Description: i18n.T("onboarding.file.bug_report.desc"),
			ScaffoldFills: []string{
				i18n.T("onboarding.file.bug_report.scaffold.1"),
				i18n.T("onboarding.file.bug_report.scaffold.2"),
				i18n.T("onboarding.file.bug_report.scaffold.3"),
			},
			AgentFills: []string{},
		},
		{
			Path: ".github/ISSUE_TEMPLATE/feature_request.md", FillLevel: FilledByScaffold,
			Description: i18n.T("onboarding.file.feature_request.desc"),
			ScaffoldFills: []string{
				i18n.T("onboarding.file.feature_request.scaffold.1"),
				i18n.T("onboarding.file.feature_request.scaffold.2"),
				i18n.T("onboarding.file.feature_request.scaffold.3"),
			},
			AgentFills: []string{},
		},

		// --- Quality ---
		{
			Path: ".pre-commit-config.yaml", FillLevel: FilledByScaffold,
			Description: i18n.T("onboarding.file.precommit.desc"),
			ScaffoldFills: []string{
				i18n.T("onboarding.file.precommit.scaffold.1"),
				i18n.T("onboarding.file.precommit.scaffold.2"),
				i18n.T("onboarding.file.precommit.scaffold.3"),
			},
			AgentFills: []string{
				i18n.T("onboarding.file.precommit.agent.1"),
				i18n.T("onboarding.file.precommit.agent.2"),
			},
		},
		{
			Path: "CONTRIBUTING.md", FillLevel: FilledByScaffold,
			Description: i18n.T("onboarding.file.contributing.desc"),
			ScaffoldFills: []string{
				i18n.T("onboarding.file.contributing.scaffold.1"),
				i18n.T("onboarding.file.contributing.scaffold.2"),
				i18n.T("onboarding.file.contributing.scaffold.3"),
				i18n.T("onboarding.file.contributing.scaffold.4"),
			},
			AgentFills: []string{
				i18n.T("onboarding.file.contributing.agent.1"),
				i18n.T("onboarding.file.contributing.agent.2"),
			},
		},
		{
			Path: "SECURITY.md", FillLevel: FilledPartial,
			Description: i18n.T("onboarding.file.security.desc"),
			ScaffoldFills: []string{
				i18n.T("onboarding.file.security.scaffold.1"),
				i18n.T("onboarding.file.security.scaffold.2"),
				i18n.T("onboarding.file.security.scaffold.3"),
			},
			AgentFills: []string{
				i18n.T("onboarding.file.security.agent.1"),
				i18n.T("onboarding.file.security.agent.2"),
			},
		},

		// --- Scripts ---
		{
			Path: "scripts/setup.sh", FillLevel: FilledByScaffold,
			Description: i18n.T("onboarding.file.setup.desc"),
			ScaffoldFills: []string{
				i18n.T("onboarding.file.setup.scaffold.1"),
				i18n.T("onboarding.file.setup.scaffold.2"),
				i18n.T("onboarding.file.setup.scaffold.3"),
				i18n.T("onboarding.file.setup.scaffold.4"),
			},
			AgentFills: []string{
				i18n.T("onboarding.file.setup.agent.1"),
				i18n.T("onboarding.file.setup.agent.2"),
				i18n.T("onboarding.file.setup.agent.3"),
			},
		},
	}

	FileMap = make(map[string]MDFile, len(Files))
	for _, f := range Files {
		FileMap[f.Path] = f
	}
}
