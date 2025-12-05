// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React, {useState} from 'react';
import {FormattedMessage} from 'react-intl';
import {Modal} from 'react-bootstrap';

interface Props {
    show: boolean;
    onHide: () => void;
    onConfirm: (startTime: number, endTime: number, level: string) => void;
}

const DateRangeModal: React.FC<Props> = ({show, onHide, onConfirm}) => {
    const [startDate, setStartDate] = useState('');
    const [endDate, setEndDate] = useState('');
    const [level, setLevel] = useState('standard');

    const handleConfirm = () => {
        let startTime: number;
        let endTime: number;

        if (startDate) {
            startTime = new Date(startDate).getTime();
        } else {
            // Default to 24 hours ago
            startTime = Date.now() - (24 * 60 * 60 * 1000);
        }

        if (endDate) {
            endTime = new Date(endDate).getTime();
        } else {
            endTime = Date.now();
        }

        onConfirm(startTime, endTime, level);
        onHide();
    };

    return (
        <Modal
            show={show}
            onHide={onHide}
            backdrop='static'
        >
            <Modal.Header closeButton={true}>
                <Modal.Title>
                    <FormattedMessage
                        id='ai.summary.dateRangeModal.title'
                        defaultMessage='Summarize Channel'
                    />
                </Modal.Title>
            </Modal.Header>
            <Modal.Body>
                <div className='form-group'>
                    <label htmlFor='startDate'>
                        <FormattedMessage
                            id='ai.summary.dateRangeModal.startDate'
                            defaultMessage='Start Date'
                        />
                    </label>
                    <input
                        id='startDate'
                        type='datetime-local'
                        className='form-control'
                        value={startDate}
                        onChange={(e) => setStartDate(e.target.value)}
                    />
                    <small className='form-text text-muted'>
                        <FormattedMessage
                            id='ai.summary.dateRangeModal.startDateHelp'
                            defaultMessage='Leave empty for 24 hours ago'
                        />
                    </small>
                </div>

                <div className='form-group'>
                    <label htmlFor='endDate'>
                        <FormattedMessage
                            id='ai.summary.dateRangeModal.endDate'
                            defaultMessage='End Date'
                        />
                    </label>
                    <input
                        id='endDate'
                        type='datetime-local'
                        className='form-control'
                        value={endDate}
                        onChange={(e) => setEndDate(e.target.value)}
                    />
                    <small className='form-text text-muted'>
                        <FormattedMessage
                            id='ai.summary.dateRangeModal.endDateHelp'
                            defaultMessage='Leave empty for now'
                        />
                    </small>
                </div>

                <div className='form-group'>
                    <label htmlFor='summaryLevel'>
                        <FormattedMessage
                            id='ai.summary.dateRangeModal.level'
                            defaultMessage='Summary Level'
                        />
                    </label>
                    <select
                        id='summaryLevel'
                        className='form-control'
                        value={level}
                        onChange={(e) => setLevel(e.target.value)}
                    >
                        <option value='brief'>
                            Brief (2-3 sentences)
                        </option>
                        <option value='standard'>
                            Standard (Key points and decisions)
                        </option>
                        <option value='detailed'>
                            Detailed (Comprehensive summary)
                        </option>
                    </select>
                </div>
            </Modal.Body>
            <Modal.Footer>
                <button
                    className='btn btn-tertiary'
                    onClick={onHide}
                >
                    <FormattedMessage
                        id='ai.summary.dateRangeModal.cancel'
                        defaultMessage='Cancel'
                    />
                </button>
                <button
                    className='btn btn-primary'
                    onClick={handleConfirm}
                >
                    <FormattedMessage
                        id='ai.summary.dateRangeModal.generate'
                        defaultMessage='Generate Summary'
                    />
                </button>
            </Modal.Footer>
        </Modal>
    );
};

export default DateRangeModal;

