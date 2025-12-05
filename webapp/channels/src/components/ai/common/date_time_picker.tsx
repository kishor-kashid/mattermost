// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React from 'react';

interface Props {
    value: Date | null;
    onChange: (date: Date | null) => void;
}

const DateTimePicker: React.FC<Props> = ({value, onChange}) => {
    const formatDateTimeLocal = (date: Date | null): string => {
        if (!date) return '';
        
        const year = date.getFullYear();
        const month = String(date.getMonth() + 1).padStart(2, '0');
        const day = String(date.getDate()).padStart(2, '0');
        const hours = String(date.getHours()).padStart(2, '0');
        const minutes = String(date.getMinutes()).padStart(2, '0');
        
        return `${year}-${month}-${day}T${hours}:${minutes}`;
    };

    const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const dateValue = e.target.value;
        if (dateValue) {
            onChange(new Date(dateValue));
        } else {
            onChange(null);
        }
    };

    return (
        <input
            type='datetime-local'
            className='form-control'
            value={formatDateTimeLocal(value)}
            onChange={handleChange}
        />
    );
};

export default DateTimePicker;

