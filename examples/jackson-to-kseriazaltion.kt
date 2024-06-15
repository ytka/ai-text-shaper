package example

import kotlinx.serialization.SerialName
import kotlinx.serialization.Serializable
import kotlinx.serialization.Transient

@Serializable
open class Paper(
    @SerialName("serial") val serial: Long,
    @SerialName("id") val id: String,
    @SerialName("title") val title: String,
    @SerialName("description") val description: String
)

@Serializable
data class Example(
    val serial: Long,
    val id: String,
    val title: String,
    @Transient
    val description: String = ""
)