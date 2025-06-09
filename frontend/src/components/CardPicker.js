import React, { useState } from 'react';
import style from "./CardPicker.module.css";

const CardPicker = ({selectedCard, availableCards = [], allowCustom = false, selectCallback }) => {
    const [customCard, setCustomCard] = useState('');

    const handleCardSelect = (card) => {
        selectCallback(card);
    };

    const handleCustomCardChange = (event) => {
        setCustomCard(event.target.value);
    };

    const handleCustomCardSubmit = () => {
        if (customCard) {
            selectCallback(customCard);
        }
    };

    return (
        <div className="row justify-content-center mt-5">
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
                        type="number" maxLength="3" pattern="\d{3}"  
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
