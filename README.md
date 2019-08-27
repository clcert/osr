# Observatorio de Seguridad de la Red (OSR)

El Observatorio de Seguridad de la Red busca entregar información actualizada del estado de seguridad de la red chilena, a través de la agregación de datos de distintos proveedores.

## Módulos

1. **Tasks**: Maneja las tareas del sistema
1. **Models**: Crea y actualiza las tablas de modelo de datos.
1. **Remote**: Administra y ejecuta comandos en servidores relaconados y usados por el OSR.
1. **Scheduler**: Se encarga de agendar y coordinar las aciiones de mantenimiento e importación de datos del sistema. (Pendiente)
1. **Health**: Monitorea la salud de los componentes de OSR. (Pendiente)

## Cómo compilar

Habiendo instalado Go v1.12.7 o superior, ejecutar estos comandos: 
```
   git clone https://github.com/clcert/osr
   cd osr
   dep ensure
   go install
  
```
El comando "osr" estará disponible globalmente.

## Otras dependencias:
* Base de datos Postgresql
* libpcap-dev en ubuntu, o libpcap-devel en yum (para el importer de darknet)

## Prerrequisitos del sistema

Se recomienda instalar OSR en su propio usuario, de forma de aislar su comportamiento del de otros en el sistema.

## Inicializar OSR
Ejecutar este comando para crear los usuarios de la base de datos:
```
   osr initdb
```
Este comando creará 2 usuarios en la base de datos Postgresql: osr_writer y osr_reader. En caso que los usuarios existan, los borrará y los creará nuevamente.

## OSR Importer

Módulo encargado de importar información de escaneos activos y pasivos del CLCERT.

### Crear modelos de bases de datos
```
   osr models createdb
```
Cada vez que se cree un nuevo tipo de datos, hay que ejecutar este comando para crear la base de datos respectiva.

### Importers existentes:

* Importer de dominios nuevos y eliminados NIC.cl
* Importer de Nombres y números de ASN según CIDR Report.
* Importer de datos de escaneo de DNS
* Importer de datos de la Darknet del CLCERT
* Importer de datos de categorías de dominios definidos por el CLCERT
* Importer de datos de subredes asociadas a países y de subredes asociadas a ASNs de Maxmind Geolite2
* Importer de subredes chilenas según RIPE
* Importer de dominios en el ranking alexa top 1M
