import React from 'react';

import '../../styles/main.css';

type Props = {
    label?: string;
};

export const LoadingSpinner = ({label}: Props) => (
    <div className='ai-suite-loading'>
        <div className='ai-suite-loading__spinner'/>
        {label && <span className='ai-suite-loading__label'>{label}</span>}
    </div>
);

export default LoadingSpinner;

