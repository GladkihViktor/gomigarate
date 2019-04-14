create or replace function check_migrate(vFilename migrations.filename%type) 
returns migrations.id%type
as 
$$
declare 
	vId migrations.id%type; 
begin 
/* Выполняем проверку миграции.
 * Если миграци нет то вернем  0 
 * иначе ID миграции
 */

	select m.id
	into vId
	from migrations m 
	where upper(trim(m.filename)) = upper(trim(vFilename));
	
	 if vId is null then 
	 	vId:=0;
	 end if;
  return vId;
end;
$$ language plpgsql; 


create or replace function insert_migration(vFilename migrations.filename%type)
returns migrations.id%type
as 
$$
declare vId  migrations.id%type;
begin 
	/*Вставка информации о миграции*/
	insert into migrations(filename, applied)
	values(vFilename, true)
	returning id into vId;
	
	return vId;
	exception when others then 
		     raise  'Error: %', SQLERRM;
end;
$$ language plpgsql;