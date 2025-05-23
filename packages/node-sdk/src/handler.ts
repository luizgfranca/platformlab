import { PromptType, ToolEventDto, EventTypeValue, ToolProperties } from '@synthreon/core'
import { EventEmitter } from 'node:events'
import { Execution } from './execution'
import { ToolComponents } from './tool-components'

export type HandlerExtraOptions = {
    name?: string;
    description?: string;
}

export type ToolHandlerDefinition = {
    id: string,
    function: ToolFunction,
    extraOptions?: HandlerExtraOptions
}

export type ToolFunction = (kit: ToolComponents) => Promise<string>;

type Status = 'announcing' | 'waiting_ack' | 'connected' | 'disconnected';

// TODO: reorganize this later to make this unaccessible to users
export class Handler {
    #status: Status;
    #definition: ToolHandlerDefinition;
    #bus: EventEmitter;
    #executions: Execution[]

    #handlerId?: string;
    #announcementId?: string

    constructor(
        definition: ToolHandlerDefinition,
        bus: EventEmitter
    ) {
        this.#definition = definition
        this.#bus = bus
        this.#status = 'announcing'
        this.#executions = []
    }

    start() {
        console.debug('starting handler')
        this.#performAnnouncement()
    }

    #handleEvent(event: ToolEventDto) {
        console.debug('[handler] handling event', event)
        if (this.#status !== 'connected') {
            console.warn('DROPPING event becuase handler is in disconnected state')
        }

        switch (event.type) {
            case EventTypeValue.InteractionOpen:
                this.#onOpenEvent(event)
                return
            default:
                console.error('DROPPING: invalid event type for handler: ', event)
        }
    }

    #onOpenEvent(event: ToolEventDto) {
        const execution = new Execution(
            this.#bus,
            this.#definition,
            (event) => this.#forwardExecutionEvent(event)
        )

        // FIXME: infinitely growing, devise a way to trim it down
        this.#executions.push(execution)
        execution.start(event);
    }

    #sendToBus(event: ToolEventDto) {
        this.#bus.emit('send', event)
    }

    #forwardExecutionEvent(event: ToolEventDto) {
        event.handler_id = this.#handlerId;
        event.tool = this.#definition.id;
        this.#sendToBus(event)
    }

    #performAnnouncement() {
        console.debug('sending announcement')
        this.#sendToBus(
            this.#getAnnouncementEvent()
        )

        this.#status = 'waiting_ack'
        this.#bus.on(
            `announcement/${this.#definition.id}`,
            (event: ToolEventDto) => this.#onAnnouncementResponse(event)
        )
    }

    #onAnnouncementResponse(event: ToolEventDto) {
        console.debug("handling announcement response");
        if (!event || !event.type) {
            console.error('invalid announcement response event format', event)
            this.#status = 'disconnected'
            return;
        }

        switch (event.type) {
            case EventTypeValue.AnnouncementACK:
                if (!event.handler_id) {
                    console.error('invalid ack response event format', event)
                    this.#status = 'disconnected'
                    return;
                }

                console.info(`handler for ${this.#definition.id} successfully registered`);
                this.#handlerId = event.handler_id;
                this.#announcementId = event.announcement_id;
                this.#status = 'connected';

                console.debug(`registering handler ${this.#handlerId} to handle events`)
                this.#bus.on(this.#handlerId, (event: ToolEventDto) => this.#handleEvent(event))

                return
            case EventTypeValue.AnnouncementNACK:
                console.error('received NACK trying to register handler with reason:', event.reason)
                process.exit(1);

            default:
                console.error('invalid event type was received:', event)
                this.#status = 'disconnected'
                return;
        }
    }

    #getAnnouncementEvent(): ToolEventDto {
        const properties: ToolProperties | undefined = this.#definition.extraOptions
            ? this.#definition.extraOptions
            : undefined

        return {
            type: EventTypeValue.AnnouncementHandler,
            tool: this.#definition.id,
            tool_properties: properties
        }
    }
}
