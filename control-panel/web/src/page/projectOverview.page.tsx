import { ProjectHeader } from "@/component/projectHeader";
import { ProjectSidebar } from "@/component/projectSidebar";
import { useProvider } from "@/context/root";
import { ToolDto } from "@/dto/tool.dto";
import AuthService from "@/service/auth.service";
import ToolView from "@/view/toolView";
import { useCallback, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";

const ProjectOverviewPage = () => {
    const [selectedTool, setSelectedTool] = useState<string | undefined>();

    const provider = useProvider();
    const navigate = useNavigate();
    const params = useParams();

    const projectAcronym = params.projectAcronym;
    const project = provider.getProjects()
        .find(project => project.acronym === projectAcronym);

    const onToolSelection = useCallback((tool: ToolDto) => {
        console.log(`tool ${tool.name} selected from project ${project?.name}`)
        setSelectedTool(tool.acronym);
    }, [])

    const onLogout = useCallback(() => {
        console.log('onLogout')
        AuthService.logout();
        // TODO: create navigator that already computes prefix path
        navigate(`${import.meta.env.PL_PATH_PREFIX}/login`)
    }, [])

    const key = `${project?.acronym}:${selectedTool}`;
    return projectAcronym ? (
        <div className="flex flex-col min-h-screen">
            <ProjectHeader projectName={project?.name ?? ''} onLogoutClick={onLogout}/>
            <div className="grid grid-cols-5 flex-1 text-zinc-100">
                <div className="col-span-1 flex flex-1">
                    <ProjectSidebar projectAcronym={projectAcronym} onSelect={onToolSelection}/>
                </div>
                <main className="col-span-4 flex-1">
                    {project && selectedTool ? <ToolView key={key} project={project.acronym} tool={selectedTool} /> : <></>}
                </main>
            </div>
        </div>
    ) : <div>Internal error.</div>
    // FIXME: find a better way to handle this kind of internal error
}

export { ProjectOverviewPage }
