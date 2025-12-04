import React from 'react';

type Props = {
    summary: string;
};

type Block =
    | {type: 'text'; content: string}
    | {type: 'bullets'; items: string[]};

export const SummaryContent = ({summary}: Props) => {
    const blocks = React.useMemo(() => parseSummary(summary), [summary]);

    return (
        <div className='ai-suite-summary__content'>
            {blocks.map((block, index) => {
                if (block.type === 'text') {
                    return (
                        <p
                            key={`text-${index}`}
                            className='ai-suite-summary__paragraph'
                        >
                            {block.content}
                        </p>
                    );
                }

                return (
                    <ul
                        key={`list-${index}`}
                        className='ai-suite-summary__list'
                    >
                        {block.items.map((item, itemIndex) => (
                            <li key={`item-${index}-${itemIndex}`}>
                                {item}
                            </li>
                        ))}
                    </ul>
                );
            })}
        </div>
    );
};

const parseSummary = (content: string): Block[] => {
    const lines = content.split('\n');
    const blocks: Block[] = [];

    let currentText = '';
    let currentBullets: string[] = [];

    const flushText = () => {
        if (currentText.trim()) {
            blocks.push({type: 'text', content: currentText.trim()});
        }
        currentText = '';
    };

    const flushBullets = () => {
        if (currentBullets.length > 0) {
            blocks.push({type: 'bullets', items: currentBullets});
        }
        currentBullets = [];
    };

    for (const line of lines) {
        const trimmed = line.trim();
        if (!trimmed) {
            flushText();
            flushBullets();
            continue;
        }

        if (/^[-*•]/.test(trimmed)) {
            flushText();
            currentBullets.push(trimmed.replace(/^[-*•]\s?/, '').trim());
            continue;
        }

        flushBullets();
        currentText += `${trimmed} `;
    }

    flushText();
    flushBullets();
    return blocks;
};

export default SummaryContent;


