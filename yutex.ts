// Register a new language
monaco.languages.register({ id: "yutex" });

// Register a tokens provider for the language
monaco.languages.setMonarchTokensProvider("yutex", {
    defaultToken: '',
    tokenPostfix: '.yutex',

    tokenizer: {
        root: [
            // section
            [/^(\\section\{)([^{}]*)(})(.*)$/, ['tag.id.pug', '', 'tag.id.pug', 'comment']],
            // code
            [/^(\\begin\{code}\{)([^{}]*)(})(.*)$/, {
                cases: {
                    '$4': ['tag.id.pug', '', 'tag.id.pug', { token: 'comment', next: '@codeBlock', nextEmbedded: '$2' }],
                    '': ['tag.id.pug', '', { token: 'tag.id.pug', next: '@codeBlock', nextEmbedded: '$2' }, 'comment'],
                },
            }],
            // math
            [/^(\\begin\{math})(.*)$/, {
                cases: {
                    '$2': ['tag.id.pug', { token: 'comment', next: '@mathBlock' }],
                    '': [{ token: 'tag.id.pug', next: '@mathBlock' }, 'comment'],
                },
            }],
            // html
            [/^(\\begin\{html})(.*)$/, {
                cases: {
                    '$2': ['tag.id.pug', { token: 'comment', next: '@htmlBlock', nextEmbedded: 'text/html' }],
                    '': [{ token: 'tag.id.pug', next: '@htmlBlock', nextEmbedded: 'text/html' }, 'comment'],
                },
            }],
            // table
            [/^(\\begin\{table})(.*)$/, {
                cases: {
                    '$2': ['tag.id.pug', { token: 'comment', next: '@table' }],
                    '': [{ token: 'tag.id.pug', next: '@table' }, 'comment'],
                },
            }],
            // sample
            [/^(\\begin\{sample}\{)([1-9]\d*)(})(.*)$/, {
                cases: {
                    '$4': ['tag.id.pug', 'number', 'tag.id.pug', { token: 'comment', next: '@sample' }],
                    '': ['tag.id.pug', 'number', { token: 'tag.id.pug', next: '@sample' }, 'comment'],
                },
            }],
            // mixcode
            [/^(\\begin\{mixcode})(.*)$/, {
                cases: {
                    '$2': ['tag.id.pug', { token: 'comment', next: '@mixcode' }],
                    '': [{ token: 'tag.id.pug', next: '@mixcode' }, 'comment'],
                },
            }],
            // paragraph
            [/^(\\begin\{paragraph})(.*)$/, {
                cases: {
                    '$2': ['tag.id.pug', { token: 'comment', next: '@paragraph' }],
                    '': [{ token: 'tag.id.pug', next: '@paragraph' }, 'comment'],
                },
            }],
            // blockquote
            [/^(\\begin\{blockquote})(.*)$/, {
                cases: {
                    '$2': ['tag.id.pug', { token: 'comment', next: '@blockquote' }],
                    '': [{ token: 'tag.id.pug', next: '@blockquote' }, 'comment'],
                },
            }],
            // others
            [/^.*$/, 'comment'],
        ],
        inline: [
            [/(\\link\{)([^{}]*)(}\{)([^{}]*)(})/, ['variable', '', 'variable', '', 'variable']],
            [/(\\text\{)([^{}]*)(}\{)([^{}]*)(})/, ['string', '', 'string', '', 'string']],
            [/\\(?:newline|space)/, 'number'],
            [/(\\math\{)([^{}]*)(})/, ['number.hex', '', 'number.hex']],
            [/(\\math\[)([^\[\]]*)(])/, ['number.hex', '', 'number.hex']],
            [/(\\html\{)([^{}]*)(})/, ['key', '', 'key']],
            [/(\\html\[)([^\[\]]*)(])/, ['key', '', 'key']],
        ],
        codeBlock: [
            [/^\\end\{code}\s*$/, { token: 'tag.id.pug', next: '@pop', nextEmbedded: '@pop' }],
        ],
        mathBlock: [
            [/^(\\end\{math})(.*)$/, [{ token: 'tag.id.pug', next: '@pop' }, 'comment']],
        ],
        htmlBlock: [
            [/^\\end\{html}\s*$/, { token: 'tag.id.pug', next: '@pop', nextEmbedded: '@pop' }],
        ],
        table: [
            [/^(\\end\{table})(.*)$/, [{ token: 'tag.id.pug', next: '@pop' }, 'comment']],
        ],
        sample: [
            [/^\\sample\{input}\s*$/, { token: 'tag.id.pug', next: '@sampleInput' }],
            [/^.*$/, 'comment'],
        ],
        sampleInput: [
            [/^\\sample\{output}\s*$/, { token: 'tag.id.pug', next: '@sampleOutput' }],
        ],
        sampleOutput: [
            [/^\\end\{sample}\s*$/, { token: 'tag.id.pug', next: '@popall' }],
        ],
        mixcode: [
            [/^(\\end\{mixcode})(.*)$/, [{ token: 'tag.id.pug', next: '@pop' }, 'comment']],
            [/^(\\begin\{code}\{)([^{}]*)(})(.*)$/, {
                cases: {
                    '$4': ['tag.id.pug', '', 'tag.id.pug', { token: 'comment', next: '@codeBlock', nextEmbedded: '$2' }],
                    '': ['tag.id.pug', '', { token: 'tag.id.pug', next: '@codeBlock', nextEmbedded: '$2' }, 'comment'],
                },
            }],
            [/^.*$/, 'comment'],
        ],
        paragraph: [
            [/^(\\end\{paragraph})(.*)$/, [{ token: 'tag.id.pug', next: '@pop' }, 'comment']],
            { include: '@inline' },
        ],
        blockquote: [
            [/^(\\end\{blockquote})(.*)$/, [{ token: 'tag.id.pug', next: '@pop' }, 'comment']],
            { include: '@inline' },
        ],
    }
});
