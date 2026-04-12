import { useState, useEffect } from "react";
import { Link } from 'react-router-dom';
import { getCategoriesForSet, getSets, searchSets } from "../api";
import { type DjSet, type Category } from "../types"

function Library() {
    const [djSets, setDjSets] = useState<DjSet[]>([]);
    const [searchQuery, setSearchQuery] = useState("");

    useEffect(() => {
        if (searchQuery === "") {
            getSets().then(data => setDjSets(data));
        } else {
            searchSets(searchQuery).then(data => setDjSets(data))
        }
    }, [searchQuery]);

    const [setCategories, setSetCategories] = useState<Record<string, Category[]>>({});

    const groupedSets = djSets.reduce((groups, set) => {
        if (!groups[set.DjName]) {
            groups[set.DjName] = [];
        }
        groups[set.DjName].push(set);
        return groups;
    }, {} as Record<string, DjSet[]>);

    useEffect(() => {
        if (djSets.length === 0) return;
        
        djSets.forEach(async (set) => {
            const cats = await getCategoriesForSet(set.ID);
            setSetCategories(prev => ({ ...prev, [set.ID]: cats || [] }));
        });
    }, [djSets]);

    return (
        <div className="Library">
            <div style={{ display: "flex", gap: "10px", marginBottom: "20px"}}>
                <Link to="/add">Add New Set</Link>
                <Link to="/categories">Add New Category</Link>
                <input type="text" placeholder="Search Sets" onChange={(e) => setSearchQuery(e.target.value)}/>
            </div>            
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
                                {setCategories[set.ID]?.map((cat) => (
                                    <span key={cat.ID}>{cat.Name} </span>
                                ))}
                            </Link>

                        </div>
                    )}
                </div>
            ))}
        </div>
    )
}

export default Library;