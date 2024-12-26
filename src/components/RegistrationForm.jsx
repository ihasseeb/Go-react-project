import React, { useState } from "react";
import "../App.css"; // Link your CSS file here
import axios from "axios";

const RegistrationForm = () => {
  const [formData, setFormData] = useState({
    firstName: "",
    lastName: "",
    email: "",
    phone: "",
    dob: "",
    address: "",
  });

  const [message, setMessage] = useState("");

  const handleChange = (e) => {
    setFormData({ ...formData, [e.target.name]: e.target.value });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();

    const backendURL = "http://127.0.0.1:8000/api/users"; // Backend API URL

    try {
      const response = await axios.post(backendURL, formData);

      if (response.status === 201 || response.status === 200) {
        setMessage("✔️ Data sent successfully!");
      } else {
        setMessage("❌ Something went wrong.");
      }
    } catch (error) {
      setMessage(
        `❌ Failed to send data. ${
          error.response?.data?.message || error.message
        }`
      );
    }
  };

  return (
    <div className="Registration-Form">
      <h1 className="title">Registration Form</h1>
      <form className="contact-form" onSubmit={handleSubmit}>
        {message && (
          <p style={{ color: message.startsWith("✔️") ? "green" : "red" }}>
            {message}
          </p>
        )}
        <div className="form-field">
          <input
            type="text"
            id="first-name"
            name="firstName"
            className={`input-text ${formData.firstName ? "not-empty" : ""}`}
            value={formData.firstName}
            onChange={handleChange}
            required
          />
          <label htmlFor="first-name" className="label">
            First Name
          </label>
        </div>
        <div className="form-field">
          <input
            type="text"
            id="last-name"
            name="lastName"
            className={`input-text ${formData.lastName ? "not-empty" : ""}`}
            value={formData.lastName}
            onChange={handleChange}
            required
          />
          <label htmlFor="last-name" className="label">
            Last Name
          </label>
        </div>
        <div className="form-field">
          <input
            type="email"
            id="email"
            name="email"
            className={`input-text ${formData.email ? "not-empty" : ""}`}
            value={formData.email}
            onChange={handleChange}
            required
          />
          <label htmlFor="email" className="label">
            Email
          </label>
        </div>
        <div className="form-field">
          <input
            type="text"
            id="phone"
            name="phone"
            className={`input-text ${formData.phone ? "not-empty" : ""}`}
            value={formData.phone}
            onChange={handleChange}
            required
          />
          <label htmlFor="phone" className="label">
            Phone
          </label>
        </div>
        <div className="form-field">
          <input
            type="date"
            id="dob"
            name="dob"
            className={`input-text ${formData.dob ? "not-empty" : ""}`}
            value={formData.dob}
            onChange={handleChange}
            required
          />
          <label htmlFor="dob" className="label">
            Date of Birth
          </label>
        </div>
        <div className="form-field">
          <textarea
            id="address"
            name="address"
            className={`input-text ${formData.address ? "not-empty" : ""}`}
            value={formData.address}
            onChange={handleChange}
            required
          ></textarea>
          <label htmlFor="address" className="label">
            Address
          </label>
        </div>
        <button type="submit" className="submit-btn">
          Submit
        </button>
      </form>
    </div>
  );
};

export default RegistrationForm;
