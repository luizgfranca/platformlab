import { PlatformConnection } from "../src/platform-connection";

const connection = new PlatformConnection({
    endpoint: 'ws://localhost:8080/api/tool/provider/ws',
    credentials: {
        username: 'test@test.com',
        password: 'password'
    },
    toolFunction: async ({ io }) => {
        const input = await io.prompt(
            'Input ping to receive pong',
            'string'
        )

        if(input.toLowerCase() === 'ping') {
            return 'pong!'
        }

        throw '...'
    }
})

connection.listen();