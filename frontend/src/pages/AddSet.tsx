import { useState } from "react";
import { createSet } from "../api";
import { useNavigate } from 'react-router-dom'

function AddSet() {
    const navigate = useNavigate();

    const [formData, setFormData] = useState ({
        url: "",
        djName: ""
    })
    
    const handleInput = (e: React.ChangeEvent<HTMLInputElement>) => {
        setFormData({...formData, [e.target.name]: e.target.value});
    }

    const handleSubmit = async () => {
        await createSet(formData.url, formData.djName);
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
            <button type="button" onClick={handleSubmit}>Add Set</button>
        </form>
    )
}
export default AddSet;