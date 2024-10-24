import { useCallback, useMemo, useState } from "react";
import { Result } from "./result";
import { Prompt } from "./prompt";

type DisplayElement = {
    type: string;
    label: string;
    text: string;
    description: string;
    name: string;
};

type DisplayResult = {
    success: boolean;
    message: string;
};

type DisplayPrompt = {
    title: string;
    type: string;
}

type DisplayDefinitionType = 'result' | 'view' | 'prompt' | string

export type DisplayDefinition = {
    type: DisplayDefinitionType;
    elements?: DisplayElement[];
    result?: DisplayResult;
    prompt?: DisplayPrompt;
};

export type DsiplayRendererProps = {
    definition: DisplayDefinition
}

export function DisplayRenderer(props: DsiplayRendererProps) {
    switch(props.definition.type) {
        case 'result':
            return (
                <Result success={props.definition.result?.success ?? false}>{props.definition.result?.message ?? ''}</Result>
            )
        case 'prompt':
            return (
                <Prompt title={props.definition.prompt?.title ?? ''} onSubmit={(value) => console.log('prompt submit', value)}/>
            )
    }

}