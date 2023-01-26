package polyxia.faas.gateway

import org.springframework.boot.autoconfigure.SpringBootApplication
import org.springframework.boot.runApplication
import org.springframework.cloud.gateway.route.RouteLocator
import org.springframework.cloud.gateway.route.builder.PredicateSpec
import org.springframework.cloud.gateway.route.builder.RouteLocatorBuilder
import org.springframework.context.annotation.Bean

@SpringBootApplication
class GatewayApplication {
    @Bean
    fun myRoutes(builder: RouteLocatorBuilder): RouteLocator? {
        return builder.routes()
                .route {
                    it
                            .path("/create")
                            .filters {
                                it
                                        .rewritePath("/(?<segment>.*)", "/api/v0/workloads.create")
                            }
                            .uri("http://localhost:10000")
                }
                .route {
                    it
                            .path("/invoke")
                            .filters {
                                it
                                        .rewritePath("/(?<segment>.*)", "/api/v0/instances.create")
                            }
                            .uri("http://localhost:10000")
                }.build()
    }
}

fun main(args: Array<String>) {
    runApplication<GatewayApplication>(*args)
}
