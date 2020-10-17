package main

import (
	"context"
	"log"
	"time"
	"math/rand"
	//"fmt"
	// Llamamos el paquete de gRPC
	"google.golang.org/grpc"
	// Llamamos el compilado que nos generó protoc
	pb "./logistica"
)

const (
	address     = "localhost:50054" // Definimos por que host y puerto nos comunicamos
)

type Paquete struct {
	Id string
	Tipo int32
	Valor int32
	Origen string
	Destino string
	Intentos int
	FechaEntrega string
}

var paqCam1 []Paquete //paquetes que han estado o estan en el camion 1
var paqCam2 []Paquete
var paqCam3 []Paquete
var abortar1 bool
var abortar2 bool

func entregaExitosa() bool {
	rand.Seed(time.Now().UnixNano())
	randomI := rand.Intn(10) 
	if randomI < 8 {
		log.Println(randomI)
		return true		
	}
	log.Println(randomI)
	return false
}

func camion1(){
	abortar1 = false
	var enCamion1 []string
	// configurar la conexión al servidor
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Error en la conexión: %v", err)
	}
	defer conn.Close()
	c := pb.NewPaqueteClient(conn)
	abortar := 0 
	for len(enCamion1) < 2 {
		//time.Sleep(time.Second)
				// Enviamos la petición al servidor
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		r, err := c.SolicitudPaquete(ctx, &pb.RepartoRequest{IdCamion: 1})

		if err != nil {
			log.Fatalf("Error en la petición: %v desde camion1()", err)
		}
		// Imprimimos la respuesta con el numero de seguimiento generado
		log.Println("Enviando solicitud de paquetes a entregar")

		if len(r.Id) > 0{
			nuevo := true
			for i := 0; i < len(paqCam1); i++ {
				if paqCam1[i].Id == r.Id{
					paqCam1[i].Intentos += 1
					nuevo = false
					break
				}
			}
			if nuevo{
				paqCam1 = append(paqCam1, Paquete{
					Id: r.Id, 
					Tipo: r.Tipo, 
					Valor: r.Valor, 
					Origen: r.Origen, 
					Destino: r.Destino, 
					Intentos: 1,
					FechaEntrega: ""})
					log.Println("Se recibio el paquete: ", r.Id)
			}
			enCamion1 = append(enCamion1, r.Id)
		}else{
			if abortar < 1{
				abortar += 1
				time.Sleep(time.Second) // si no hay paquetes para repartir en la central, esperamos un segundo y volvemos a consultar
				log.Println("No se recibio paquetes. Se intentara denuevo en 1[s]")				
			}else{
				break
			}

		}
		
	}
	if abortar != 10 || len(enCamion1) != 0{
		if len(enCamion1) > 0{
			log.Println("Los id-paquete en camion son: ", enCamion1)	
		}
		
		// se informa el exito en la entrega de los dos paquetes
		for i := 0; i < len(enCamion1); i++ {
			exitoEntrega := entregaExitosa()
			//log.Println("Exito en la entrega:", exitoEntrega )
			if exitoEntrega{
				for j := 0; j < len(paqCam1); j++ {
					if paqCam1[j].Id == enCamion1[i]{
						paqCam1[j].FechaEntrega = time.Now().Format("2006-01-02 15:04:05")
					}
				}
			}
			// Enviamos la petición al servidor
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			r, err := c.InformeReparto(ctx, &pb.InformeRequest{Id: enCamion1[i], Entregado: exitoEntrega})

			if err != nil {
				log.Fatalf("Error en la petición: %v ", err)
			}
			log.Println("Enviando informe de entrega...")
			log.Println("Reporte enviado con exito:", r.Recibido)		
		}
	}
}

func camion2(){
	abortar2 = false
	var enCamion2 []string
	// configurar la conexión al servidor
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Error en la conexión: %v", err)
	}
	defer conn.Close()
	c := pb.NewPaqueteClient(conn)
	abortar := 0 
	for len(enCamion2) < 2 {
		//time.Sleep(time.Second)
				// Enviamos la petición al servidor
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		r, err := c.SolicitudPaquete(ctx, &pb.RepartoRequest{IdCamion: 2})

		if err != nil {
			log.Fatalf("Error en la petición: %v desde Camion2", err)
		}
		// Imprimimos la respuesta con el numero de seguimiento generado
		log.Println("Enviando solicitud de paquetes a entregar")

		if len(r.Id) > 0{
			nuevo := true
			for i := 0; i < len(paqCam2); i++ {
				if paqCam2[i].Id == r.Id{
					paqCam2[i].Intentos += 1
					nuevo = false
					break
				}
			}
			if nuevo{
				paqCam2 = append(paqCam2, Paquete{
					Id: r.Id, 
					Tipo: r.Tipo, 
					Valor: r.Valor, 
					Origen: r.Origen, 
					Destino: r.Destino, 
					Intentos: 1,
					FechaEntrega: ""})
					log.Println("Se recibio el paquete: ", r.Id)
			}
			enCamion2 = append(enCamion2, r.Id)
		}else{
			if abortar < 1{
				abortar += 1
				time.Sleep(time.Second) // si no hay paquetes para repartir en la central, esperamos un segundo y volvemos a consultar
				log.Println("No se recibio paquetes. Se intentara denuevo en 1[s]")				
			}else{
				break
			}

		}
		
	}
	if abortar != 10 || len(enCamion2) != 0{
		if len(enCamion2) > 0 {
			log.Println("Los id-paquete en camion son: ", enCamion2)	
		}
		
		// se informa el exito en la entrega de los dos paquetes
		for i := 0; i < len(enCamion2); i++ {
			exitoEntrega := entregaExitosa()
			//log.Println("Exito en la entrega:", exitoEntrega )
			if exitoEntrega{
				for j := 0; j < len(paqCam2); j++ {
					if paqCam2[j].Id == enCamion2[i]{
						paqCam2[j].FechaEntrega = time.Now().Format("2006-01-02 15:04:05")
					}
				}
			}
			// Enviamos la petición al servidor
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			r, err := c.InformeReparto(ctx, &pb.InformeRequest{Id: enCamion2[i], Entregado: exitoEntrega})

			if err != nil {
				log.Fatalf("Error en la petición: %v ", err)
			}
			log.Println("Enviando informe de entrega...")
			log.Println("Reporte enviado con exito:", r.Recibido)		
		}
	}
}

func camion3(){
	var enCamion3 []string
	// configurar la conexión al servidor
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Error en la conexión: %v", err)
	}
	defer conn.Close()
	c := pb.NewPaqueteClient(conn)
	abortar := 0 
	for len(enCamion3) < 2 {
		//time.Sleep(time.Second)
				// Enviamos la petición al servidor
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		r, err := c.SolicitudPaquete(ctx, &pb.RepartoRequest{IdCamion: 3})

		if err != nil {
			log.Fatalf("Error en la petición: %v desde Camion2", err)
		}
		// Imprimimos la respuesta con el numero de seguimiento generado
		log.Println("Enviando solicitud de paquetes a entregar")
		if len(r.Id) > 0{
			nuevo := true
			for i := 0; i < len(paqCam3); i++ {
				if paqCam3[i].Id == r.Id{
					paqCam3[i].Intentos += 1
					nuevo = false
					break
				}
			}
			if nuevo{
				paqCam3 = append(paqCam3, Paquete{
					Id: r.Id, 
					Tipo: r.Tipo, 
					Valor: r.Valor, 
					Origen: r.Origen, 
					Destino: r.Destino, 
					Intentos: 1,
					FechaEntrega: ""})
					log.Println("Se recibio el paquete: ", r.Id)
			}
			enCamion3 = append(enCamion3, r.Id)
		}else{
			if abortar < 1{
				abortar += 1
				time.Sleep(time.Second) // si no hay paquetes para repartir en la central, esperamos un segundo y volvemos a consultar
				log.Println("No se recibio paquetes. Se intentara denuevo en 1[s]")				
			}else{
				break
			}

		}
		
	}
	if abortar != 10 || len(enCamion3) != 0{
		if len(enCamion3) > 0{
			log.Println("Los id-paquete en camion son: ", enCamion3)	
		}
		
		// se informa el exito en la entrega de los dos paquetes
		for i := 0; i < len(enCamion3); i++ {
			exitoEntrega := entregaExitosa()
			//log.Println("Exito en la entrega:", exitoEntrega )
			if exitoEntrega{
				for j := 0; j < len(paqCam3); j++ {
					if paqCam3[j].Id == enCamion3[i]{
						paqCam3[j].FechaEntrega = time.Now().Format("2006-01-02 15:04:05")
					}
				}
			}
			// Enviamos la petición al servidor
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			r, err := c.InformeReparto(ctx, &pb.InformeRequest{Id: enCamion3[i], Entregado: exitoEntrega})

			if err != nil {
				log.Fatalf("Error en la petición: %v ", err)
			}
			log.Println("Enviando informe de entrega...")
			log.Println("Reporte enviado con exito:", r.Recibido)		
		}
	}
}



func main() {
	start := time.Now()//comienza la simulacion
	
	for{
		camion1()	
		camion2()	
		camion3()
		if time.Since(start).String ()[:2] == "1m"{
			log.Println("A pasado mas de 1[m] de simulacion, por tanto se termina la jornada de repartos")	
			break
		}

	}
	//fmt.Scanln()

	
	
	
	
}