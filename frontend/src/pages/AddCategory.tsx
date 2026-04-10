import { useState } from "react";
import { createCategory } from "../api";
import { useNavigate } from 'react-router-dom'

function AddCategory() {
    const navigate = useNavigate();

    const [formData, setFormData] = useState ({
        category: ""
    })
    
    const handleInput = (e: React.ChangeEvent<HTMLInputElement>) => {
        setFormData({...formData, [e.target.name]: e.target.value});
    }

    const handleSubmit = async () => {
        await createCategory(formData.category);
        navigate('/');
    }

    return (
        <form>
            <div>
                <label>Category</label>
                <input type="text" name="category" onChange={handleInput}/>
            </div>
            <button type="button" onClick={handleSubmit}>Add Category</button>
        </form>
    )
}
export default AddCategory;