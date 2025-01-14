import { DisplayDefinition, DisplayRenderer, DsiplayRendererProps, Field } from "@/component/displayRenderer"
import { Prompt } from "@/component/prompt";
import { useCallback, useMemo, useState } from "react"

// const event = {
//     "class": "operation",
//     "type": "display",
//     "project": "proj-x",
//     "tool": "tool-y",
//     "display": {
//         "type": "result",
//         "result": {
//             "success": true,
//             "message": "Hello user input",
//         }
//     }
// }

// const testEventPrompt = {
//     "class": "operation",
//     "type": "display",
//     "project": "proj-x",
//     "tool": "tool-y",
//     "display": {
//         "type": "prompt",
//         "prompt": {
//             "title": "Add some text:",
//             "type": "string",
//         }
//     }
// }


type ToolEvent = {
    "class": string;
    "type": string;
    "project": string;
    "tool": string;
    "display": DisplayDefinition;
}

export function ToolViewExperimentsPage() {
    const [event, setEvent] = useState<ToolEvent | null>(null)

    console.log('e', event)
    
    const ws = useMemo(() => {
        const ws = new WebSocket(`${import.meta.env.PL_BACKEND_URL}/api/tool/client/ws`)
        ws.addEventListener('open', () => {
            console.log('socket open')

            ws.send(JSON.stringify({
                "class": "interaction",
                "type": "open",
                "project": "proj-x",
                "tool": "tool-p",
            }))
        })

        ws.addEventListener('message', (e) => {
            console.log(`recv: ${e.data}`)
            setEvent(JSON.parse(e.data))
        })

        return ws
    }, [])
 
    const sendInputInteraction = (fields: Field[]) => {
        ws.send(JSON.stringify({
            "class": "interaction",
            "type": "input",
            "project": "proj-x",
            "tool": "tool-p",
            input: {
                fields
            }
        }))
    }

    if(!event) {
        return (
            <div className="bg-zinc-900 text-zinc-100 h-screen">
                    <div className="container mx-auto px-4 py-8">
                    <h1 className="text-3xl font-bold mb-6">Tool Sandbox</h1>

                    {/* <div className="flex justify-center">
                        <div className="w-4/5">
                            <Prompt title="Test prompt component" onSubmit={(value) => console.log('onsubmit', value)} />
                        </div>
                    </div>
                     */}
                </div>
            </div>
        )       
    }
    return (
        <div className="bg-zinc-900 text-zinc-100 h-screen">
            <div className="container mx-auto px-4 py-8">
                <h1 className="text-3xl font-bold mb-6">Tool Sandbox</h1>
                <DisplayRenderer definition={event.display} onSumission={(fields) => sendInputInteraction(fields)}/>
            </div>
        </div>
    )
}