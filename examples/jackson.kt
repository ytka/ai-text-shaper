package example
import com.fasterxml.jackson.annotation.*

open class Paper(
    @JsonProperty("serial") val serial: Long,
    @JsonProperty("id") val id: String,
    @JsonProperty("title") val title: String,
    @JsonProperty("description") val description: String
)

data class Example(
    val serial: Long,
    val id: String,
    val title: String,
    @JsonIgnore
    val description: String,
)
