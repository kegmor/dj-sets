import { useState, useEffect } from "react";
import { Link } from 'react-router-dom';
import { getSets } from "../api";
import { type DjSet } from "../types"

function Library() {
    const [djSets, setDjSets] = useState<DjSet[]>([]);

    useEffect(() => {
        getSets().then(data => setDjSets(data));
    }, [])

    return (
        <div className="Library">
            <Link to="/add">Add New Set</Link>
            {djSets.map((set: DjSet) => (
                <Link to={`sets/${set.ID}`} key={set.ID}>
                    <h2>{set.DjName}</h2>
                    <p>{set.Title}</p>
                    <p>{set.ChannelName}</p>
                </Link>
            ))}
        </div>
    );
}

export default Library;