import React from 'react';

type Props = {
    points?: string[];
};

export const KeyPoints = ({points}: Props) => {
    if (!points || points.length === 0) {
        return null;
    }

    return (
        <ul className='ai-link-summary__points'>
            {points.map((p, idx) => (
                <li
                    key={idx}
                    className='ai-link-summary__point'
                >
                    {p}
                </li>
            ))}
        </ul>
    );
};

