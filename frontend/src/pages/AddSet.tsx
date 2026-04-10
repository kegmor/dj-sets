import { useEffect, useState } from "react";
import { createSet, getCategories, addCategoryToSet } from "../api";
import { useNavigate } from 'react-router-dom'
import type { Category } from "../types";

function AddSet() {
    const navigate = useNavigate();

    const [djSetData, setDjSetData] = useState ({
        url: "",
        djName: ""
    })

    const [categories, setCategoriesData] = useState <Category[]>([]);
    const [selectedCategories, setSelectedCategories] = useState <string[]>([]);

    const handleCheckBox = (name: string) => {
        if (selectedCategories.includes(name)) {
            setSelectedCategories(selectedCategories.filter(c => c !== name));
        } else {
            setSelectedCategories([...selectedCategories, name]);
        }
    }

    useEffect (() => {
        getCategories().then(data => setCategoriesData(data))
    }, [])
    
    const handleInput = (e: React.ChangeEvent<HTMLInputElement>) => {
        setDjSetData({...djSetData, [e.target.name]: e.target.value});
    }

    const handleSubmit = async () => {
        const createdSet = await createSet(djSetData.url, djSetData.djName);
        for (const name of selectedCategories) {
            await addCategoryToSet(createdSet.ID, name);
        }
        navigate('/');
    }

    return (
        <form>
            <div>
                <label>URL</label>
                <input type="text" name="url" onChange={handleInput}/>
            </div>
            <div>
                <label>DJ Name</label>
                <input type="text" name="djName" onChange={handleInput}/>
            </div>
            <div>
                <label>Music Category</label>
                {categories.map((cat) => (
                    <label key={cat.ID}>
                    <input type="checkbox" onChange={() => handleCheckBox(cat.Name)} />
                    {cat.Name}
                    </label>
                ))}
                
            </div>
            <button type="button" onClick={handleSubmit}>Add Set</button>
        </form>
    )
}
export default AddSet;
