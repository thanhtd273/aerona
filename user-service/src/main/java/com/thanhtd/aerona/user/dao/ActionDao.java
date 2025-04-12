package com.thanhtd.aerona.user.dao;

import com.thanhtd.aerona.user.model.Action;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.query.Param;
import org.springframework.stereotype.Repository;

import java.util.List;

@Repository
public interface ActionDao extends JpaRepository<Action, String> {

    @Query("SELECT u FROM Action u WHERE u.status <> 0")
    List<Action> findAll();

    @Query("SELECT u FROM Action u WHERE u.status <> 0 AND u.actionId = :actionId")
    Action findByActionId(@Param(value = "actionId") String actionId);

    @Query("SELECT u FROM Action u WHERE u.status <> 0 AND u.code = :code")
    Action findByCode(@Param(value = "code") Integer code);
}
