import { TextField } from "@mui/material"
import { FC } from "react"
import { useFormContext } from "react-hook-form"

interface IReactHookFormTextFieldProps {
    label: string
    name: string
    value?: string
}

const ReactHookFormTextField: FC<IReactHookFormTextFieldProps> = ({ label, name, value }: IReactHookFormTextFieldProps) => {
    const {
        register,
        formState: { errors },
    } = useFormContext()

    return (
        <TextField
            label={label}
            variant="outlined"
            value={value}
            error={!!errors[name]}
            helperText={errors[name]?.message ?? ''}
            fullWidth
            margin="dense"
            {...register(name)}
        />
    )
}

export default ReactHookFormTextField
