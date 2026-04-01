import { useState, useEffect } from "react";
import { getSets } from "../api";
import { type DjSet } from "../types"

function Library() {
    const [djSets, setDjSets] = useState<DjSet[]>([]);

    useEffect(() => {
        getSets().then(data => setDjSets(data));
    }, [])

    return (
        <div className="Library">
            {djSets.map((set: DjSet) => (
                <div key={set.ID}>
                    <h2>{set.DjName}</h2>
                    <p>{set.Title}</p>
                    <p>{set.ChannelName}</p>
                </div>
            ))}
        </div>
    );
}

export default Library;