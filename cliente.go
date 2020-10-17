package main

import (
	"context"
	"log"
	"time"
	"fmt"
	"io"
	"os"
	"strconv"
	"encoding/csv"
	// Llamamos el paquete de gRPC
	"google.golang.org/grpc"
	// Llamamos el compilado que nos generó protoc
	pb "./logistica"
)

const (
	address     = "localhost:50053" // Definimos por que host y puerto nos comunicamos
)

type Paquete struct {
	Id string
	Tipo int
	Nombre string
	Valor int
	Origen string
	Destino string
	Seguimiento string
}

func seguimientoPaq(paquetes []Paquete){
	start := time.Now()//comienza la simulacion
	time.Sleep(time.Second)
	// configurar la conexión al servidor
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Error en la conexión: %v", err)
	}
	defer conn.Close()
	c := pb.NewPaqueteClient(conn)

	var entregados [40]int
	for i := 0; i < len(paquetes); i++ {
		entregados[i] = 0
	}
	f:= true
	for f{
		time.Sleep(2*time.Second)
		for i := 0; i < len(paquetes); i++ {

			if entregados[i] == 0{//no queremos mostrar nuevamente el estado de los paquetes que ya fueron entregados

				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				r, err := c.SeguimientoPaquete(ctx, &pb.SeguimientoRequest{
					NumSeguimiento: paquetes[i].Seguimiento})
				if err != nil {
					log.Fatalf("Error en la petición: %v", err)
				}
				// Imprimimos la respuesta con el estado
				if len(r.Estado) > 2{
					log.Println("El estado del paquete es: ", r.Estado)	
					if r.Estado == "Recibido"{
						entregados[i] = 1 

					}else if r.Estado == "No Recibido"{
						entregados[i] = 1 						
					}
				}
			}
			if time.Since(start).String ()[:2] == "1m"{
				f = false
				log.Fatal("A pasado mas de 1[m] de simulacion, por tanto se termina la ejecucion")	
				//log.Println(time.Since(start))
				break
			}

		}		
	}
}


func enviarPaq(paquetes []Paquete){
	time.Sleep(2*time.Second)

		// configurar la conexión al servidor
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Error en la conexión: %v", err)
	}
	defer conn.Close()
	c := pb.NewPaqueteClient(conn)
	//time.Sleep(time.Second)
	var it = 0
	for it != len(paquetes){
	// Enviamos la petición al servidor
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.EnviarPaquete(ctx, &pb.PaqueteRequest{
		IdPaquete: paquetes[it].Id,
		Tipo : int32(paquetes[it].Tipo), //0=es retail, 1=pyme, 2=prioritario
		Nombre : paquetes[it].Nombre,//nombre producto
		Valor: int32(paquetes[it].Valor),
		Origen: paquetes[it].Origen,
		Destino:  paquetes[it].Destino})
	if err != nil {
		log.Fatalf("Error en la petición: %v desde enviarPaq()", err)
	}
	// Imprimimos la respuesta con el numero de seguimiento generado
	log.Println("Enviando el paquete %d ", it)
	log.Println("Codigo de seguimiento para el paquete: ", r.NumSeguimiento)
	paquetes[it].Seguimiento = r.NumSeguimiento
	it += 1
	}
}

func clienteRetail() {
	var paquetes []Paquete //estructura que contendra todos los paquetes
	log.Println("Ha iniciado la simulacion como Retail")	
	//se definen los slices que almacenaran las columnas del csv
	var id []string
	var producto []string
	var temp []string
	var tienda []string
	var destino []string
	csvfile, err := os.Open("./input/retail.csv") // abrir el csv
	if err != nil {
		log.Fatalln("El archivo .csv no pudo ser abierto", err)
	}
	r := csv.NewReader(csvfile) // Parse the file

	
	for { // Iteratar en los registrosdel csv
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		//se asigna cada columna del csv a un slice
		id = append(id, record[0])
		producto = append(producto, record[1])
		temp = append(temp, record[2])
		tienda =  append(tienda, record[3])
		destino = append(destino, record[4])
	}

	//se quita el nombre de la columna de cada slice
	id = id[1:]
	producto = producto[1:]
	temp = temp[1:]
	tienda = tienda[1:]
	destino = destino[1:]

	//se parsea los valores leidos como string desde el csv a entero
	valor := make([]int, len(temp))
    for i := range temp {
        valor[i], _ = strconv.Atoi(temp[i])
    }

	for i := 0; i < len(valor); i++ {
		paquetes = append(paquetes, Paquete{
			Id: id[i], 
			Tipo: 0, 
			Nombre: producto[i], 
			Valor: valor[i], 
			Origen: tienda[i], 
			Destino: destino[i], 
			Seguimiento: ""})
		
	}

	go enviarPaq(paquetes)
	time.Sleep(time.Second)
	go seguimientoPaq(paquetes)
}

func clientePyme() {
	//estructura que contendra todos los paquetes
	var paquetes []Paquete

	log.Println("Ha iniciado la simulacion como Pyme")
	
	//se definen los slices que almacenaran las columnas del csv
	var id []string
	var producto []string
	var temp []string
	var tienda []string
	var destino []string
	var prioritario []string

	// abrir el csv
	csvfile, err := os.Open("./input/pymes.csv")
	if err != nil {
		log.Fatalln("El archivo .csv no pudo ser abierto", err)
	}
	// Parse the file
	r := csv.NewReader(csvfile)
	//r := csv.NewReader(bufio.NewReader(csvfile))

	// Iterate through the records
	for {
		// Read each record from csv
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		//se asigna cada columna del csv a un slice
		id = append(id, record[0])
		producto = append(producto, record[1])
		temp = append(temp, record[2])
		tienda =  append(tienda, record[3])
		destino = append(destino, record[4])
		prioritario = append(prioritario, record[5])
	}

	//se quita el nombre de la columna de cada slice
	id = id[1:]
	producto = producto[1:]
	temp = temp[1:]
	tienda = tienda[1:]
	destino = destino[1:]
	prioritario = prioritario[1:]

	//se parsea los valores leidos como string desde el csv a entero
	valor := make([]int, len(temp))
    for i := range temp {
        valor[i], _ = strconv.Atoi(temp[i])
    }


	// se lleva los leido desde el csv a la estructura que contendra los paquetes
	for i := 0; i < len(valor); i++ {
		if prioritario[i]=="1"{
			paquetes = append(paquetes, Paquete{
				Id: id[i], 
				Tipo: 2, 
				Nombre: producto[i], 
				Valor: valor[i], 
				Origen: tienda[i], 
				Destino: destino[i], 
				Seguimiento: ""})
		}else{
			paquetes = append(paquetes, Paquete{
				Id: id[i], 
				Tipo: 1, 
				Nombre: producto[i], 
				Valor: valor[i], 
				Origen: tienda[i], 
				Destino: destino[i], 
				Seguimiento: ""})
		}
	}

	go enviarPaq(paquetes)
	time.Sleep(time.Second)
	go seguimientoPaq(paquetes)
}

func main() {
	//start := time.Now()//comienza la simulacion
	//clientePyme()
	//time.Sleep(time.Second)
	go clienteRetail()
	//time.Sleep(2 * time.Second)
	go clientePyme()
	//time.Sleep(2 * time.Second)
	go clienteRetail()
	//time.Sleep(2 * time.Second)


	fmt.Println("A pasado mas de 1[m] de simulacion, por tanto se termina la ejecucion")
	fmt.Scanln()


}
