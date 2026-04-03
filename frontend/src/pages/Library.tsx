import { useState, useEffect } from "react";
import { Link } from 'react-router-dom';
import { getSets } from "../api";
import { type DjSet } from "../types"

function Library() {
    const [djSets, setDjSets] = useState<DjSet[]>([]);

    useEffect(() => {
        getSets().then(data => setDjSets(data));
    }, [])

    const groupedSets = djSets.reduce((groups, set) => {
        if (!groups[set.DjName]) {
            groups[set.DjName] = [];
        }
        groups[set.DjName].push(set);
        return groups;
    }, {} as Record<string, DjSet[]>);

    return (
        <div className="Library">
            <Link to="/add">Add New Set</Link>
            {Object.entries(groupedSets).map(([djName, sets]) => (
                <div key={djName}>
                    <h2>{djName}</h2>
                    {sets.map((set: DjSet) => 
                        <div key={set.ID}>
                            <iframe
                                width="560"
                                height="315"
                                src={`https://www.youtube.com/embed/${set.VideoID}`}
                                title={set.Title}
                                allowFullScreen
                            />
                            <Link to={`sets/${set.ID}`} key={set.ID}>
                                <p>{set.Title}</p>
                                <p>{set.ChannelName}</p>
                            </Link>
                        </div>
                    )}
                </div>
            ))}
        </div>
    )
}

export default Library;