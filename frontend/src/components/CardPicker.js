import React, { useState } from 'react';
import style from "./CardPicker.module.css";

const CardPicker = ({selectedCard, availableCards = [], allowCustom = false, selectCallback }) => {
    const [customCard, setCustomCard] = useState('');

    const handleCardSelect = (card) => {
        selectCallback(card);
    };

    const handleCustomCardChange = (event) => {
        const value = event.target.value;

        if (value === '' || (value >= 1 && value <= 999)) {
            setCustomCard(value);
        } else {
            setCustomCard(1);
        }
    };

    const handleCustomCardSubmit = () => {
        if (customCard) {
            selectCallback(customCard);
        }
    };

    return (
        <div className="row justify-content-center mt-5">
            <div 
                className={`${style.card} ${selectedCard === -1 ? style.selected : ''}`}
                onClick={() => handleCardSelect(-1)}
            >
                <i className="bi bi-eye-fill"></i>
            </div>
            {availableCards.map((value, index) => {
                return (
                    <div 
                        key={index} 
                        className={`${style.card} ${selectedCard === value ? style.selected : ''}`}
                        onClick={() => handleCardSelect(value)}
                    >
                        {value}
                    </div>
                );
            })}
            {allowCustom && (
                <div className={`${style.cardCustom} ${selectedCard === customCard ? style.selected : ''}`}>
                    <input
                        className={style.cardInput} 
                        type="number" 
                        min="1" 
                        max="999" 
                        value={customCard} 
                        onChange={handleCustomCardChange} 
                        placeholder='???'
                    /><br/>
                    <button className='btn btn-primary' onClick={handleCustomCardSubmit}>Ok</button>
                </div>
            )}
        </div>
    );
};

export default CardPicker;
