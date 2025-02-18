package providermodule

import (
	"log"
	commonmodule "platformlab/controlpanel/modules/common"
	projectmodule "platformlab/controlpanel/modules/project"
	toolmodule "platformlab/controlpanel/modules/tool"
	"platformlab/controlpanel/modules/toolentity"
	tooleventmodule "platformlab/controlpanel/modules/toolevent"
)

type Orchestrator interface {
	ForwardEventToClient(e *tooleventmodule.ToolEvent)
}

// FIXME: add project and tool deregistration handling
// TODO: create abstractions for managerServices

// FIXME: when a provider disconnects it should notify the orchestrator in order for
// it to cancel the context
type ProviderManagerService struct {
	orchestrator Orchestrator

	projectService *projectmodule.ProjectService

	contextProviderResolver        ContextProviderResolver
	projectAndToolProviderResolver ProjectAndToolProviderResolver

	// TODO: is this list irrelevant?
	providers []*Provider
}

// OnProviderDisconnection implements Manager.
func (p *ProviderManagerService) OnProviderDisconnection(provider *Provider) {
	if provider == nil {
		log.Fatalln("provider to be unregistered from providermanager should not be null")
	}

	p.log("provider ", provider.ID, " disconnected, degeristering mappings")
	p.projectAndToolProviderResolver.UnregisterProviderEntries(provider)
	p.contextProviderResolver.UnregisterProviderEntries(provider)

	// TODO: could this have a race condition if the provider unregisters right after connecting?
	for i, it := range p.providers {
		if it.ID == provider.ID {
			commonmodule.RemoveFromUnorderedSlice(p.providers, i)
			break // provider should not have duplicates in the list
		}
	}
}

func NewProviderManagerService(
	orchestrator Orchestrator,
	projectService *projectmodule.ProjectService,
	toolService *toolmodule.ToolService,
) *ProviderManagerService {
	return &ProviderManagerService{
		orchestrator:   orchestrator,
		projectService: projectService,

		contextProviderResolver:        ContextProviderResolver{},
		projectAndToolProviderResolver: ProjectAndToolProviderResolver{},
		providers:                      []*Provider{},
	}
}

// DistributeEvent implements Manager.
// called by the managed providers
func (p *ProviderManagerService) DistributeEvent(e *tooleventmodule.ToolEvent) {
	if e.Type == tooleventmodule.EventTypeCommandFinish {
		p.contextProviderResolver.Unregister(e.ContextId)
	}

	e.HandshakeId = ""
	e.ExecutionId = ""

	p.orchestrator.ForwardEventToClient(e)
}

// RegisterProviderProjectAndTool implements Manager.
func (p *ProviderManagerService) RegisterProviderProjectAndTool(m *ProviderToolMapping) {
	p.projectAndToolProviderResolver.Register(m.Project.Acronym, m.Tool.Acronym, m.Provider)
}

// FindProject implements Manager.
func (p *ProviderManagerService) FindProject(acronym string) (*projectmodule.Project, error) {
	maybeProject, err := p.projectService.FindByAcronym(acronym)
	if err != nil {
		return nil, err
	}

	if maybeProject == nil {
		return nil, &commonmodule.GenericLogicError{Message: "project not fonud"}
	}

	return maybeProject, nil
}

// FindTool implements Manager.
func (p *ProviderManagerService) FindTool(project *projectmodule.Project, acronym string) (*toolmodule.Tool, error) {
	p.log(
		"looking for announced tool: \n",
		"project: ", project.Acronym, "\n",
		"tool: ", acronym,
	)

	maybeTool, err := p.projectService.FindToolByAcronym(project, acronym)
	if err != nil {
		return nil, err
	}

	if maybeTool == nil {
		return nil, &commonmodule.GenericLogicError{Message: "tool from project not fonud"}
	}

	p.log("tool found:", maybeTool)
	return maybeTool, nil
}

func (p *ProviderManagerService) EntityConnection(entity toolentity.ToolEntityAdapter) {
	p.providers = append(p.providers, NewProvider(p, entity))
}

func (p *ProviderManagerService) SendEvent(e *tooleventmodule.ToolEvent) error {
	p.log("event received by providerManager")
	if e.ContextId == "" || e.Tool == "" || e.Project == "" {
		log.Fatalln("[ProviderManagerService] unexptected event attributes when reaching providerManager", e)
	}

	p.log("trying to route")
	err := p.contextProviderResolver.TryRouteEvent(e)
	if err != nil {
		switch err.(type) {
		case *ContextNotFounError:
			p.log("context not found, creating")
			provider, err := p.projectAndToolProviderResolver.Resolve(e.Project, e.Tool)
			if err != nil {
				p.log("provider resolution error: ", err.Error())
				return err
			}

			p.log("registering new context")
			p.contextProviderResolver.Register(e.ContextId, provider)

			p.log("routing event to provider")
			err = p.contextProviderResolver.TryRouteEvent(e)
			if err != nil {
				log.Fatalln("unexpected behavior: could not route event even after context setup")
			}

			return nil
		case *commonmodule.GenericLogicError:
			return err
		}
	}

	return nil
}

func (s *ProviderManagerService) UnregisterContext(ctxid string) {
	s.log("unregistering context: ", ctxid)
	provider := s.contextProviderResolver.TryResolve(ctxid)
	if provider == nil {
		log.Fatalln("[ProviderManagerService] unexpected: trying to unregister unexisting context: ", ctxid)
	}

	provider.UnregisterContext(ctxid)
	s.contextProviderResolver.Unregister(ctxid)
}

func (p *ProviderManagerService) log(v ...any) {
	x := append([]any{"[ProviderManagerService]"}, v...)

	log.Println(x...)
}

// TODO: I dont know if i like creating this struct just for this,
// maybe there's a better approach
type ProviderToolMapping struct {
	Provider *Provider
	Project  *projectmodule.Project
	Tool     *toolmodule.Tool
}
