# cli-naive-replication: an (unsuccessful) attempt to keep in sync three redis databases
`cli-naive-replication` es implementado para evidenciar el problema de inconsistencias que surge en un ambiente
distribuido compuesto por al menos 2 clientes y un grupo de nodos de almacenamiento. 

El problema es que no se puede garantizar que los datos en diferentes nodos de bases de datos estén
siempre en sincronía. En sistemas distribuidos, este problema de inconsistencias es muy común sin embargo es difícil
de resolver y requiere la inclusión de un protocolo State Machine Replication (SMR) el cual no es aplicado en 
esta implementación. 
El objetivo de esta implementación es comprender de forma experimental el problema de inconsistencias 
que surge en un intento de replicación que no incluye un protocolo SMR. 
Así mismo, evaluamos el rendimiento (throughput y latencia) de las operaciones de escritura y lectura al servicio de 
almacenamiento compuesto por 3 nodos.


