# ✈️ Aerona Flight Booking System

Aerona is a microservices-based flight booking system designed to handle high-performance search, secure booking, and reliable payment processing. The system leverages modern technologies and architecture principles to ensure scalability, maintainability, and a seamless user experience.

---

## 🚀 Technologies Used

- **Backend**: Combination of **Golang** and **Java Spring Boot**
- **Architecture**: **Event-driven architecture (choreography pattern)** using **Apache Kafka** for inter-service communication
- **Search & Caching**: **Elasticsearch** and **Redis** for high-speed search and temporary storage
- **Data Storage**:
  - **PostgreSQL** for reliable and consistent storage of critical data such as payments and bookings
- **CI/CD**: **GitLab CI** for automating the deployment of services

---

## 🧱 Microservices Overview

The system is composed of 7 core services:

1. **user-service** _(Spring Boot)_  
   Handles user authentication, authorization, and user management.

2. **airline-integration-service** _(Golang)_  
   Integrates with external airline providers via API to fetch flight information and store it in Elasticsearch. Also supports webhooks for real-time synchronization of flight data.

3. **flight-search-service** _(Golang)_  
   Provides APIs for searching and filtering flights, checking seat availability, and saving search history.

4. **booking-service** _(Spring Boot)_  
   Manages the full booking flow—from selecting a flight to confirming the reservation.

5. **payment-service** _(Spring Boot)_  
   Integrates with **PayPal** to securely process payments.

6. **ticket-service** _(Golang)_  
   Generates e-tickets and invoices, and uploads them to **AWS S3** for storage.

7. **notification-service** _(Golang)_  
   Sends emails and push notifications to users regarding their bookings and updates.

---

## 🧩 System Architecture

![Aerona System Architecture](./aerona-infra/docs/aerona-architecture.png)

---

## 🛠️ Work in Progress

- [ ] Deploying the system to **Kubernetes**
- [ ] Hosting the system on **AWS**

---

## 📫 Contact

For questions or contributions, feel free to reach out or create an issue. Contributions are welcome!
