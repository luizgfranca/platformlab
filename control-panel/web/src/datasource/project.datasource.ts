import { ProjectDto } from "@/dto/project.dto";
import { Datasource } from "@/lib/datasource";
import ProjectService from "@/service/project.service";

export class ProjectDatasource extends Datasource<ProjectDto[]> {
    async fetch(): Promise<ProjectDto[]>{
        return ProjectService.queryProjects();
    }
}