import React, {ReactNode} from 'react';

import '../../styles/main.css';

type Props = {
    title: string;
    onClose: () => void;
    children: ReactNode;
    footer?: ReactNode;
};

export const Modal = ({title, onClose, children, footer}: Props) => (
    <div className='ai-suite-modal__backdrop'>
        <div className='ai-suite-modal'>
            <div className='ai-suite-modal__header'>
                <h3>{title}</h3>
                <button
                    type='button'
                    className='ai-suite-modal__close'
                    aria-label='Close'
                    onClick={onClose}
                >
                    ×
                </button>
            </div>
            <div className='ai-suite-modal__body'>
                {children}
            </div>
            {footer && (
                <div className='ai-suite-modal__footer'>
                    {footer}
                </div>
            )}
        </div>
    </div>
);

export default Modal;

