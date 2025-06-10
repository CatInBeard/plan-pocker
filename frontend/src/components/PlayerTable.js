import React, { useState, useEffect, useRef } from 'react';
import './PlayerTable.css';

const PlayerTable = ({ voters,revealedValue, revealAction, startNewAction }) => {
    const [averageCardValue, setAverageCardValue] = useState(null);
    const [playerPositions, setPlayerPositions] = useState({});
    const tableRef = useRef(null);
    const playerWidth = 80;
    const playerHeight = 60;
    const aspectRatio = 7 / 3;

    useEffect(() => {
        const calculatePlayerPositions = () => {
            if (!tableRef.current) return;

            const tableWidth = tableRef.current.offsetWidth;
            const tableHeight = tableRef.current.offsetHeight;
            const numPlayers = voters.length;
            const positions = {};

            const radiusX = tableWidth * 0.45;
            const radiusY = tableHeight * 0.45;

            const centerX = tableWidth / 2;
            const centerY = tableHeight / 2;

            voters.forEach((voter, i) => {
                const angle = (2 * Math.PI * i) / numPlayers;

                const x = centerX + radiusX * Math.cos(angle);
                const y = centerY + radiusY * Math.sin(angle);

                const offsetX = (playerWidth / 2) * Math.cos(angle);
                const offsetY = (playerHeight / 2) * Math.sin(angle);

                positions[voter.username] = { x: x + offsetX, y: y + offsetY };
            });

            setPlayerPositions(positions);
        };

        const timeoutId = setTimeout(() => {
            calculatePlayerPositions();
        }, 0);

        window.addEventListener('resize', calculatePlayerPositions);

        return () => {
            clearTimeout(timeoutId);
            window.removeEventListener('resize', calculatePlayerPositions);
        };
    }, [voters]);

    return (
        <div className="poker-table" ref={tableRef}>
            <div className="table-center">
                {revealedValue === null && (<button className="show-all-cards-button" onClick={revealAction}>
                    Reveal
                </button> )}
                {revealedValue !== null && (<button className="start-new-game-button" onClick={startNewAction}>
                    Start new game
                </button> )}
                {revealedValue !== null && (
                    <p className="average-card-value">Average: {revealedValue}</p>
                )}
            </div>
            {voters.map((voter) => (
                <div
                    key={voter.username}
                    className={(voter.vote === "ready") ? "closed player" : "player"}
                    style={{
                        left: playerPositions[voter.username]?.x + 'px',
                        top: playerPositions[voter.username]?.y + 'px',
                        transform: 'translate(-50%, -50%)',
                    }}
                >
                    <h5 className='top-text'>{voter.username}</h5>
                    <p>
                        {!isNaN(voter.vote) ? (
                            <span>{voter.vote}</span>
                        ) : voter.vote === "waiting" ? (
                            <span>Waiting</span>
                        ) : voter.vote === "ready" ? (
                            <span>Ready</span>
                        ) : (
                            <span>{voter.vote}</span>
                        )}
                    </p>
                </div>
            ))}
        </div>
    );
};

export default PlayerTable;
