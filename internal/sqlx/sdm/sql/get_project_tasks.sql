select t_id,
       -- p_id,  --- sqlx fails if no field in model
       t_priority, t_date, t_subject from tasks where p_id =?
order by t_id