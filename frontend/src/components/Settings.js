import { useEffect, useState } from "react";
import Popup from "./Popup";

const Settings = ({ deck, allowCustomDeck, setSettings, closeSettings }) => {
    const [localDeck, setLocalDeck] = useState(deck);
    const [isCustomDeckAllowed, setIsCustomDeckAllowed] = useState(allowCustomDeck);

    useEffect(() => {
        setLocalDeck(deck.join(','));
        setIsCustomDeckAllowed(allowCustomDeck)
    }, []);

    const changeDeck = (event) => {
        setLocalDeck(event.target.value);
    };

    const toggleCustomDeck = () => {
        setIsCustomDeckAllowed(!isCustomDeckAllowed);
    };

    const saveSettings = () => {
        const numbers = localDeck.split(',')
            .map(item => item.trim())
            .map(Number)
            .filter(item => !isNaN(item));

        if (numbers.length > 0) {
            numbers.sort((a, b) => a - b);
            setSettings(numbers, isCustomDeckAllowed);
        } else {
            const powersOfTwo = Array.from({ length: 9 }, (_, i) => Math.pow(2, i));
            setSettings(powersOfTwo, isCustomDeckAllowed);
        }
    };


    return (
        <Popup header={"Game settings:"} closeAction={closeSettings}>
            <div className="form-group mb-2">
                <label>Deck:</label>
                <input
                    type="text"
                    className="form-control"
                    value={localDeck}
                    onChange={changeDeck}
                />
            </div>
            <div className="form-check">
                <input
                    type="checkbox"
                    className="form-check-input"
                    checked={isCustomDeckAllowed}
                    onChange={toggleCustomDeck}
                />
                <label className="form-check-label">Allow custom deck</label>
            </div>
            <div className="form-group">
                <button className="btn btn-primary" onClick={saveSettings}>
                    Save
                </button>
            </div>
        </Popup>
    );
};

export default Settings;
