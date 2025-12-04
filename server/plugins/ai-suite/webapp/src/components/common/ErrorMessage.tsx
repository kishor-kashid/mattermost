import React from 'react';

import '../../styles/main.css';

type Props = {
    title?: string;
    message: string;
    onRetry?: () => void;
    retryLabel?: string;
};

export const ErrorMessage = ({title = 'Something went wrong', message, onRetry, retryLabel = 'Retry'}: Props) => (
    <div className='ai-suite-error'>
        <div className='ai-suite-error__title'>{title}</div>
        <div className='ai-suite-error__body'>{message}</div>
        {onRetry && (
            <button
                className='ai-suite-error__action'
                onClick={onRetry}
                type='button'
            >
                {retryLabel}
            </button>
        )}
    </div>
);

export default ErrorMessage;

