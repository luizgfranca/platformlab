import { DisplayTypeValue, PromptType } from 'platformlab-core'
import { Execution } from './execution'
import { InputDefinition } from 'platformlab-core/tool-event/input/input.dto'


export type PromptParams = {
    title: string,
    type: PromptType
}
export type PromptFunction = (params: PromptParams) => Promise<string>
export type TextBoxFunction = (content: string) => Promise<void>

export type ToolComponents = {
    io: {
        prompt: PromptFunction
        textBox: TextBoxFunction
    }
}

function instantiatePrompt(execution: Execution): PromptFunction {
    return ({title, type}: PromptParams) =>
        new Promise((resolve, reject) => {
            console.debug('executing prompt dispatch', { title, type })
            execution.sendDisplay({
                type: DisplayTypeValue.Prompt,
                prompt: {
                    title,
                    type,
                },
            })

            // FIXME: should handle if callback is not valid anymore
            execution.onNextInput((input: InputDefinition) => {
                if (input.fields.length == 0) {
                    return reject("protocol Error: did not receive field result");
                }

                const value = input.fields[0].value
                if (!value) {
                    return reject("protocol Error: value did not arrive with input interaction");
                }

                resolve(value)
            })
        })
}

function instantiateTextBox(execution: Execution): TextBoxFunction {
    return (content: string) =>
        new Promise((resolve, _) => {
            execution.sendDisplay({
                type: DisplayTypeValue.TextBox,
                textBox: { content }
            })

            execution.onNextInput(() => {
                resolve()
            })
        })
}

function instantiateComponents(execution: Execution): ToolComponents {
    return {
        io: {
            prompt: instantiatePrompt(execution),
            textBox: instantiateTextBox(execution),
        },
    }
}

const ComponentFactory = {
    instantiateComponents,
}

export { ComponentFactory }
