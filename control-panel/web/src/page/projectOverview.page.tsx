import { ProjectHeader } from "@/component/projectHeader";
import { ProjectSidebar } from "@/component/projectSidebar";
import { useProvider } from "@/context/root";
import { ToolDto } from "@/dto/tool.dto";
import AuthService from "@/service/auth.service";
import { ToolView } from "@/view/toolView";
import { useState } from "react";
import { useNavigate } from "react-router-dom";

const defaultTools: ToolDto[] = [
    {
        id: 0,
        acronym: 'sandbox',
        name: 'Sandbox',
        description: 'Sandbox development environment for tools',
        project_id: 0
    }
]

export function ProjectOverviewPage(props: unknown) {
    const [selectedTool, setSelectedTool] = useState<string | undefined>();

    const provider = useProvider();
    const navigate = useNavigate();
    const projectAcronym = window.location.pathname.split('/')[2];

    const project = provider.getProjects().find(project => project.acronym === projectAcronym);

    const onToolSelection = (tool: ToolDto) => {
        console.log(`tool ${tool.name} selected from project ${project?.name}`)
        setSelectedTool(tool.acronym);
    }

    const onLogout = () => {
        console.log('onLogout')
        AuthService.logout();
        navigate('/login')
    }

    return (
        <div>
            <ProjectHeader projectName={project?.name ?? ''} onLogoutClick={onLogout}/>
            <div className="grid grid-cols-5 h-screen text-zinc-100">
                <div className="col-span-1">
                    <ProjectSidebar tools={defaultTools} onSelect={onToolSelection}/>
                </div>
                <main className="col-span-4">
                    <ToolView project={project?.acronym} tool={selectedTool} />
                </main>
            </div>
        </div>
    )
}