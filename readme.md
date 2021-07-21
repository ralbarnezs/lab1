# Laboratory 1

## Description
This is a laboratory of the Distributed Systems course(INF-343).

## About

 This lab aims to apply both gRPC with protocol buffers as well as RabbitMQ to simulate the operation of a distributed packet system. All this will be done in golang programming language.
 
 
# Implementacion y detalles



# Carpetas
### input :
Contiene los archivos csv con los respectivos datos de pedidos a enviar por parte de pymes y retail.

### logistica:
Contiene el compilado generado por logistica.proto en el cual se encuantran los servicios y metodos que hacen posible el enviar paquete,  realizar seguimiento al paquete y los metodos respectivos de peparto a consumirse, mediante gRPC.

# Archivos:

- **cliente.go**: Contiene los metodos necesarios para simular los requerimientos de realizar envio mediante  PrestigioExpress. Al comenzar, se puede iniciar como cliente retail o pyme y en base a esto se cargan los registros desde el respectivo csv. Luego es posible realizar el seguimiento de los paquetes enviados mediante gRPC a la central. Consideraciones  y supuestos de este archivo son:

  - Al iniciar la simulacion, se corre en conjunto dos clientes retail y un cliente pyme, mediante go goroutines para simular la concurrencia entre los clientes.

  - Una vez iniciada la simulacion con el rol respectivo, se comienza a realizar envios de paquetes, lo cual consiste en enviar los datos, uno a uno, de los registros en el csv entregado.

  - El enviar paquetes y realizar el seguimiento se realizan mediante go goroutines, para no tener que esperar a terminar una accion o interrumpirla. Ademas, el realizar seguimientos parte un segundo mas tarde para que tenga sentido.

  - El seguimiento se hace cada dos segundos, solo a los paquetes que sabemos aun no han sido recibidos.

  - Se marca los paquetes de los cuales se sabe su estado final para no volver a preguntar por ellos. Solo se pregunta por los paquetes que en un ultimo seguimiento aparecen en bodega o en reparto.

  - La simulacion termina pasado el minuto de ejecucion, por defecto.


- **reparto.go**: Contiene todo lo necesario par simular los requerimientos de pedir un paquete a la central. Luego si el paquete nunca se ha intentado entregar, se añade al registro historico del camion (un arreglo de paquetes en memoria). Sino, se busca y se incrementa el intento (pues desde finanzas estimaron que aun es posible intentarlo). En caso de que no haya paquetes en la central, se intentara recibir un paquete nuevamente en **1 segundo**. Otros supestos y requerimientos implementados son:

    - Camion 1: Primero intenta recibir un paquete de retail, si no hay paquetes de retail, intenta recibir un paquete prioritario. Una vez teniendo el paquete, espera un segundo para pedir otro. Sino hay, entonces sale a reparto.

    - Camion 2: Lo mismo que para el camion 1, solo que llega en segundo lugar(en base al apartado en el requirimiento qu hay un orden de llegada a la central, para facilitar la simulacion), y por tanto espera a que el camion 1 salga a reparto para pedir paquetes.

    - Camion 3:Llega en tercer lugo. Primero intenta recibir un paquete de prioritario, si no hay paquetes prioritarios, intenta recibir un paquete pyme. Una vez teniendo el paquete, espera un segundo para pedir otro. Sino hay, entonces sale a reparto.

    - La simulacion termina pasado el minuto de ejecucion.

  Al momento de simular el reparto, se escoge un numero al azar entre uno y diez y en base a esto se simula el exito, para dar cumplimiento a el requisito del %80 de probabilidad a ser repartido. Luego se informa a la central el exito o no en la entrega y se vuelva a pedir mas paquetes.


- **servidor.go**: Contiene implementados los metodos definidos en el compilado de gRPC ademas de las implementaciones para hacer psible la conexion con finanzas mediante RabbitMQ. Supuestos y requerimientos implementados:

  - **cliente**: Se designa un puerto en especifico para escuchar las peticiones del cliente. Al recibir una solicitud de envio, se genera un id al paquete ademas de un id de seguimiento, respondiendo al susuario con el codigo de seguimiento generado. Luego se añaden a un registro de paquetes y a una cola de paquetes dependiendo de si es envio normal, pyme o prioritario. 

  - **reparto**: Se designa un puerto diferente al de clientes. Al recibir una peticion de paquetes por parte de un camion, con su id de camion, se busca en las respectivas colas dependiendo del id del camion, y entonces se envia una respuesta con el id del paquete y demas, se quita de la cola, se incrementa en 1 los intentos, y se cambia el estado del paquete a en reparto. Luego al recibir una solicitud de informe con el id del paquete, si fue posible entregarlo, se cambia el estado del paquete a entregado, sino, se encola denuevo en su respectiva cola, segun su tipo de paquete, si finanzas determina que es rentable el entregarlo nuevamente y se cambia su estado a en bodega esperando a ser repartido,  de otra forma, su estado cambia a no recibido. Como para PrestigioExpress, el prestigio es muy importante, en el caso de que un paquete del tipo pyme no se pueda entregar, se vuelve a intentar denuevo si el numero de intentos es menor que 2 o si el costo por un segundo intento es **igual** o menor al precio del articulo. Para el caso de un paquete retail, se intenta tres veces como maximo entregarlo.


  - La simulacion termina pasado el minuto y diez segundos de ejecucion. Al terminar, se genera un archivo csv llamado **historico.csv** con el detalle final de los paquetes recibidos y gestionados durante toda la ejecucion.

