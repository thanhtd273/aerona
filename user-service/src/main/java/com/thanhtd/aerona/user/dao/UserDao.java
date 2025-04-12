package com.thanhtd.aerona.user.dao;

import com.thanhtd.aerona.user.model.User;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.query.Param;
import org.springframework.stereotype.Repository;

@Repository
public interface UserDao extends JpaRepository<User, String> {

    @Query("SELECT u FROM User u WHERE u.status <> 0 AND u.email = :email")
    User findByEmail(@Param("email") String email);

    @Query("SELECT u FROM User u WHERE u.status <> 0 AND u.userId = :userId")
    User findByUserId(@Param("userId") String userId);
}
